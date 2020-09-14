package realtimenotification

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_realtimeNotification_startPolling(t *testing.T) {
	send := make(chan []byte, 5)
	defer close(send)
	response := make(chan []byte, 5)
	defer close(response)
	ResponseFromPolling := make(chan json.RawMessage, 5)
	defer close(ResponseFromPolling)
	stopConnecting := make(chan struct{}, 5)
	defer close(stopConnecting)

	ctx, cancel := context.WithCancel(context.Background())

	login := Login{
		Authentification: Auth{
			Token: "somecredentials",
		},
		SystemOfUnits: "metric",
	}

	api := RealtimeNotificationAPI{
		ctxcancel:           cancel,
		ResponseFromPolling: ResponseFromPolling,
		stopConnecting:      stopConnecting,
		ctx:                 ctx,
		send:                send,
		login:               login,
		response:            response,
	}

	answer := `[{"data":{"realtimeAction":"CREATE","data":{"creationTime":"2020-09-10T08:04:54.747Z","deviceId":"deviceId","deviceName":"aDevice","self":"selflink","id":"1000","status":"PENDING","description":"Close relay","c8y_Relay":{"relayState":"CLOSED"}}},"channel":"/operations/2867","id":"24404966"}]`
	api.startPolling()
	api.response <- []byte(answer)
	time.Sleep(3 * time.Second)
	api.stopPolling()
	answerFromPolling := <-api.ResponseFromPolling
	expected := `{"realtimeAction":"CREATE","data":{"creationTime":"2020-09-10T08:04:54.747Z","deviceId":"deviceId","deviceName":"aDevice","self":"selflink","id":"1000","status":"PENDING","description":"Close relay","c8y_Relay":{"relayState":"CLOSED"}}}`
	assert.Equal(t, expected, string(answerFromPolling))
}
