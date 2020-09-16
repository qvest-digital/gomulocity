package realtimenotification

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	websocket "github.com/gorilla/websocket"
)

// The API does not yet surveil the Connectionsstatus.

type RealtimeNotificationAPI struct {
	timeout             time.Duration
	bufferLength        int
	ctx                 context.Context
	ctxcancel           context.CancelFunc
	login               Login
	connection          *websocket.Conn
	send                chan []byte
	response            chan []byte
	interrupt           chan os.Signal
	stopConnecting      chan struct{}
	ResponseFromPolling chan json.RawMessage
	pollingRunning      bool
	clientID            string
}

type AuthRequest struct {
	Channel                  string   `json:"channel"`
	Ext                      Login    `json:"ext"`
	Version                  string   `json:"version"`
	MinimumVersion           string   `json:"minimumVersion"`
	SupportedConnectionTypes []string `json:"supportedConnectionTypes"`
	Advice                   Advice   `json:"advice"`
}

type Login struct {
	Authentification Auth   `json:"com.cumulocity.authn"`
	SystemOfUnits    string `json:"systemOfunits"`
}

type Auth struct {
	Token     string `json:"token"`
	Tfa       string `json:"tfa,omitempty"`
	XsrfToken string `json:"xsrfToken,omitempty"`
}

type AuthResponse struct {
	Id                       string   `json:"id,omitempty"`
	Ext                      Ext      `json:"ext"`
	Channel                  string   `json:"channel"`
	Version                  string   `json:"version"`
	MinimumVersion           string   `json:"minimumVersion"`
	SupportedConnectionTypes []string `json:"supportedConnectionTypes"`
	ClientId                 string   `json:"clientId"`
	Successful               bool     `json:"successful"`
	Error                    string   `json:"error"`
}

type Ext struct {
	Ack bool `json:"ack"`
}
type SubscriptionRequest struct {
	Id           string `json:"id,omitempty"`
	ClientId     string `json:"clientId"`
	Ext          Login  `json:"ext"`
	Channel      string `json:"channel"`
	Subscription string `json:"subscription"`
}

type SubscriptionResponse struct {
	Id           string `json:"id"`
	Channel      string `json:"channel"`
	ClientId     string `json:"clientId"`
	Subscription string `json:"subscription"`
	Successful   bool   `json:"successful"`
	Error        string `json:"error"`
}

type ConnectRequest struct {
	Id             string `json:"id"`
	Ext            Login  `json:"ext"`
	Channel        string `json:"channel"`
	ClientId       string `json:"clientId"`
	ConnectionType string `json:"connectionType"`
	Advice         Advice `json:"advice"`
}

type Advice struct {
	Timeout  int `json:"timeout"`
	Interval int `json:"interval"`
}

type ConnectResponse struct {
	Id         string          `json:"id"`
	Channel    string          `json:"channel"`
	ClientId   string          `json:"clientId"`
	Successful bool            `json:"successful"`
	Data       json.RawMessage `json:"data"`
	Error      string          `json:"error"`
}

// StartRealtimeNotificationsAPI assembles all components, opens the connection via websocket. credential have to follow the pattern:"tenantid/userid:password"
func StartRealtimeNotificationsAPI(ctx context.Context, credentials, adress string, opts ...APIOption) (*RealtimeNotificationAPI, error) {
	const (
		defaultTimeout      = 5 * time.Second
		defaultBufferLength = 5
	)

	api := RealtimeNotificationAPI{
		timeout:      defaultTimeout,
		bufferLength: defaultBufferLength,
	}

	for _, opt := range opts {
		opt(api)
	}

	ctxForAPI, cancel := context.WithCancel(ctx)

	send := make(chan []byte)

	response := make(chan []byte, 10)

	stopConnecting := make(chan struct{}, 10)

	responseFromPolling := make(chan json.RawMessage, 10)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	connection, err := initConnection(ctxForAPI, adress)

	encodedCredentials := b64.StdEncoding.EncodeToString([]byte(credentials))

	login := Login{
		Authentification: Auth{
			Token: encodedCredentials,
		},
		SystemOfUnits: "metric",
	}

	api.ctxcancel = cancel
	api.ctx = ctxForAPI
	api.login = login
	api.connection = connection
	api.send = send
	api.response = response
	api.stopConnecting = stopConnecting
	api.ResponseFromPolling = responseFromPolling
	api.interrupt = interrupt

	if err != nil {
		return nil, err
	}
	api.startSendRoutine()
	api.startReadRoutine()
	api.doHandshake()
	api.startPolling()

	return &api, nil
}

//DoSubscribe subscribes to a given Channel ex "operations/1234"
func (api *RealtimeNotificationAPI) DoSubscribe(subscription string) error {
	subscriptionrequest := SubscriptionRequest{
		Channel:      "/meta/subscribe",
		Ext:          api.login,
		ClientId:     api.clientID,
		Subscription: subscription,
	}
	mustRestartPolling := false
	if api.pollingRunning {
		mustRestartPolling = true
		api.stopPolling()
	}
	out, err := json.Marshal(subscriptionrequest)
	if err != nil {
		return err
	}
	api.send <- out
	respFromSubscription := make([]SubscriptionResponse, 1)

	select {
	case <-time.After(api.timeout):
		if mustRestartPolling {
			api.startPolling()
		}
		return fmt.Errorf("timeout while waiting for subscription answer")
	case answer := <-api.response:
		json.Unmarshal(answer, &respFromSubscription)
		if !respFromSubscription[0].Successful {
			if mustRestartPolling {
				api.startPolling()
			}
			return fmt.Errorf("subscription unsuccessful with response: %v", respFromSubscription)
		}
		if mustRestartPolling {
			api.startPolling()
		}
		return nil
	}
}

// DoUnsubscribe does unsubscribe from a given channel ex "operation/1234"
func (api *RealtimeNotificationAPI) DoUnsubscribe(subscription string) error {
	subscriptionrequest := SubscriptionRequest{
		Channel:      "/meta/unsubscribe",
		Ext:          api.login,
		ClientId:     api.clientID,
		Subscription: subscription,
	}
	mustRestartPolling := false
	if api.pollingRunning {
		mustRestartPolling = true
		api.stopPolling()
	}
	out, err := json.Marshal(subscriptionrequest)
	if err != nil {
		if mustRestartPolling {
			api.startPolling()
		}
		return err
	}
	api.send <- out
	respFromSubscription := make([]SubscriptionResponse, 1)

	select {
	case <-time.After(api.timeout):
		if mustRestartPolling {
			api.startPolling()
		}
		return fmt.Errorf("timeout while waiting for unsubscription answer")
	case answer := <-api.response:
		json.Unmarshal(answer, &respFromSubscription)
		if !respFromSubscription[0].Successful {
			if mustRestartPolling {
				api.startPolling()
			}
			return fmt.Errorf("unsubscription unsuccessful with response: %v", respFromSubscription)
		}
		if mustRestartPolling {
			api.startPolling()
		}
		return nil
	}
}

func (api *RealtimeNotificationAPI) startPolling() {
	connectrequest := ConnectRequest{
		Channel:        "/meta/connect",
		Ext:            api.login,
		ClientId:       api.clientID,
		ConnectionType: "websocket",
	}
	out, _ := json.Marshal(connectrequest)
	api.pollingRunning = true
	go func() {
		for {
			select {
			case <-api.stopConnecting:
				api.pollingRunning = false
				return
			case answer := <-api.response:
				respFromConnect := []ConnectResponse{}
				json.Unmarshal(answer, &respFromConnect)
				if len(respFromConnect[0].Data) != 0 {
					api.ResponseFromPolling <- respFromConnect[0].Data
				}
			case <-time.After(500 * time.Millisecond):
				api.send <- out
			}
		}
	}()
}

func initConnection(ctx context.Context, adress string) (*websocket.Conn, error) {

	u := url.URL{Scheme: "wss", Host: adress, Path: "cep/realtime"}

	c, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (api *RealtimeNotificationAPI) startReadRoutine() {
	go func() {
		for {
			_, message, err := api.connection.ReadMessage()
			if err != nil {
			}
			api.response <- message
		}
	}()
}

func (api *RealtimeNotificationAPI) startSendRoutine() {
	go func() {
		for {
			select {
			case t := <-api.send:
				err := api.connection.WriteMessage(websocket.TextMessage, t)
				if err != nil {
					return
				}
			case <-api.interrupt:
				api.stop()
				return
			}
		}
	}()
}

func (api *RealtimeNotificationAPI) stop() {
	if api.pollingRunning {
		api.stopPolling()
	}
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.

	err := api.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		api.ctxcancel()
	}
	api.ctxcancel()
}

func (api *RealtimeNotificationAPI) doHandshake() error {
	ext := AuthRequest{
		Channel: "/meta/handshake",
		Ext:     api.login,
		Advice: Advice{
			Timeout:  60000,
			Interval: 10000,
		},
		Version:                  "1.0",
		MinimumVersion:           "1.0",
		SupportedConnectionTypes: []string{"websocket"},
	}

	out, err := json.Marshal(ext)

	if err != nil {
		return err
	}
	api.send <- out

	select {
	case <-time.After(api.timeout):
		return fmt.Errorf("timeout waiting for response from handshake")
	case answer := <-api.response:
		respFromHandshake := make([]AuthResponse, 1)
		err = json.Unmarshal(answer, &respFromHandshake)
		if err != nil {
			return err
		}
		if !respFromHandshake[0].Successful {
			return fmt.Errorf("handshake failed with response: %v", respFromHandshake)
		}
		api.clientID = respFromHandshake[0].ClientId
		return nil
	}
}

func (api *RealtimeNotificationAPI) stopPolling() {
	api.stopConnecting <- struct{}{}
	api.pollingRunning = false
}

type APIOption func(RealtimeNotificationAPI)

func withTimeout(timeout time.Duration) APIOption {
	return func(api RealtimeNotificationAPI) {
		api.timeout = timeout
	}
}

func withBufferLength(bufferlength int) APIOption {
	return func(api RealtimeNotificationAPI) {
		api.bufferLength = bufferlength
	}
}
