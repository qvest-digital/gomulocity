package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// buildingiot.test-ram.m2m.telekom.com/meta/handshake
var addr = flag.String("addr", "echo.websocket.org", "http service address")

type AuthRequest struct {
	Id                       string
	Channel                  string
	Ext                      Login
	Version                  string
	MinimumVersion           string `json:"omitempty"`
	SupportedConnectionTypes []string
}

type Login struct {
	Authentification Auth `json:"com.cumulocity.authn"`
	SystemOfUnits    string
}

type Auth struct {
	Token     string
	Tfa       string `json:"omitempty"`
	XsrfToken string `json:"omitempty"`
}

type AuthResponse struct {
	Id                       string
	Channel                  string
	Version                  string   `json:"omitempty"`
	MinimumVersion           string   `json:"omitempty"`
	SupportedConnectionTypes []string `json:"omitempty"`
	ClientId                 string   `json:"omitempty"`
	Successful               bool
	error                    string `json:"omitempty"`
}

type SubscriptiobRequest struct {
	Id           string
	Channel      string
	ClientId     string
	Subscription string
}

type SubscriptionResponse struct {
	Id           string
	Channel      string
	ClientId     string
	Subscription string
	Scuccessful  bool
	Error        string
}

type ConnectRequest struct {
	Id             string `json:"omitempty"`
	Channel        string
	ClientId       string
	ConnectionType string
	Advice         Advice `json:"omitempty"`
}

type Advice struct {
	Timeout  int `json:"omitempty"`
	Interval int `json:"omitempty"`
}

type ConnectResponse struct {
	Id         string
	Channel    string
	ClientId   string
	Successful bool
	Data       []Notification
	Error      string
}

type DisconnectRequest struct {
	Id       string `json:"omitempty"`
	Channel  string
	ClientId string
}

type DisconnectResponse struct {
	Id         string
	Channel    string
	Successful bool
	ClientId   string
	Error      string
}

var ext = AuthRequest{
	Channel: "/meta/handshake",
	Ext: Login{
		Authentification: Auth{
			Token: "login",
		},
		SystemOfUnits: "metric",
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

	response := make(chan []byte, 5)
	defer close(response)

	u := url.URL{Scheme: "ws", Host: *addr}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	log.Print(resp.Location())
	if err != nil {
		log.Fatalf("dial: %s with responsecode: %v", err, resp.StatusCode)
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
			case <-time.After(time.Second):
				return
			}
		}
	}()

	out, err := json.Marshal(ext)
	if err != nil {
		log.Fatalf("failed to Marshal: %s", err)
	}
	send <- out
	time.Sleep(3 * time.Second)
	answer := <-response
	respFromHandshake := AuthResponse{}
	json.Unmarshal(answer, respFromHandshake)
	println(respFromHandskake.ClientId)
}
