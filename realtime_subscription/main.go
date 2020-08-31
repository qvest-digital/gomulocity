package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/tarent/gomulocity/devicecontrol"

	"github.com/gorilla/websocket"
)

// buildingiot.test-ram.m2m.telekom.com/meta/handshake
var addr = flag.String("addr", "tarent-gmbh.cumulocity.com", "http service address")

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
type SubscriptiobRequest struct {
	Id           string `json:"id,omitempty"`
	Ext          Login  `json:"ext"`
	Channel      string `json:"channel"`
	ClientId     string `json:"cliendId"`
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
	Id         string                    `json:"id"`
	Channel    string                    `json:"channel"`
	ClientId   string                    `json:"clientId"`
	Successful bool                      `json:"successful"`
	Data       []devicecontrol.Operation `json:"data"`
	Error      string                    `json:"error"`
}

type DisconnectRequest struct {
	Id       string `json:"id"`
	Channel  string `json:"channel"`
	ClientId string `json:"clientId"`
}

type DisconnectResponse struct {
	Id         string `json:"id"`
	Channel    string `json:"channel"`
	Successful bool   `json:"successful"`
	ClientId   string `json:"clientId"`
	Error      string `json:"error"`
}

var credentials = b64.StdEncoding.EncodeToString([]byte(""))
var login = Login{
	Authentification: Auth{
		Token: credentials,
	},
	SystemOfUnits: "metric",
}
var ext = AuthRequest{
	Channel: "/meta/handshake",
	Ext:     login,
	Advice: Advice{
		Timeout:  60000,
		Interval: 10000,
	},
	Version:                  "1.0",
	MinimumVersion:           "1.0",
	SupportedConnectionTypes: []string{"websocket"},
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	response := make(chan []byte, 10)
	defer close(response)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "cep/realtime"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial: %s with responsecode: ", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			response <- message
			log.Printf("recv: %s", message)
		}
	}()

	send := make(chan []byte, 5)
	defer close(send)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-send:
				err := c.WriteMessage(websocket.TextMessage, t)
				log.Println("sent: ", string(t))
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-interrupt:
				log.Println("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()

	out, err := json.Marshal(ext)

	if err != nil {
		log.Fatalf("failed to Marshal: %s", err)
	}
	send <- out
	time.Sleep(1 * time.Second)
	answer := <-response
	respFromHandshake := make([]AuthResponse, 1)
	err = json.Unmarshal(answer, &respFromHandshake)
	if err != nil {
		log.Println(err)
		log.Println(string(answer))
		log.Println(respFromHandshake)
	}
	if !respFromHandshake[0].Successful {
		log.Fatal("handshake failed")
	}
	log.Println("handshake successful")
	clientID := respFromHandshake[0].ClientId
	log.Println(clientID)

	subscriptionrequest := SubscriptiobRequest{
		Channel:      "/meta/subscribe",
		Ext:          login,
		ClientId:     clientID,
		Subscription: "/operations/3329",
	}
	out, err = json.Marshal(subscriptionrequest)
	if err != nil {
		log.Fatalf("failed to Marshal: %s", err)
	}
	send <- out
	log.Println("subrequest sent")
	log.Println(string(out))
	timeout := time.After(10 * time.Second)
	log.Println("waiting for response")
	respFromSubscription := make([]SubscriptionResponse, 1)

	for {
		select {
		case <-timeout:
			break
		case answer := <-response:
			log.Println("subresponse received")
			log.Println(string(answer))
			json.Unmarshal(answer, &respFromSubscription)
			log.Println(respFromSubscription)
		}
	}
	json.Unmarshal(answer, respFromSubscription)
	if !respFromSubscription[0].Successful {
		log.Fatalf("error while subscribing: %s", respFromSubscription[0].Error)
	}
	log.Println("subscription successful")
	connectrequest := ConnectRequest{
		Channel:        "/meta/connect",
		Ext:            login,
		ClientId:       clientID,
		ConnectionType: "websocket",
	}
	out, err = json.Marshal(connectrequest)
	if err != nil {
		log.Fatalf("failed to Marshal: %s", err)
	}

	timeout = time.After(2 * time.Minute)
	for {
		select {
		case <-timeout:
			break
		case answer := <-response:

			respFromConnect := ConnectResponse{}
			json.Unmarshal(answer, respFromConnect)
			log.Println(respFromConnect)
		case <-time.After(time.Second):
			send <- out
			return
		}
	}

}
