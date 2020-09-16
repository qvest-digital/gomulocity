package realtimenotification

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

func Test_RealtimeNotification_stopAPI(t *testing.T) {
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
	con, _, err := websocket.DefaultDialer.Dial("wss://echo.websocket.org", nil)
	require.NoError(t, err)

	api := RealtimeNotificationAPI{
		connection:          con,
		ctxcancel:           cancel,
		ResponseFromPolling: ResponseFromPolling,
		stopConnecting:      stopConnecting,
		ctx:                 ctx,
		send:                send,
		login:               login,
		response:            response,
	}

	api.startPolling()
	assert.True(t, api.pollingRunning)
	api.stop()
	assert.False(t, api.pollingRunning)
	select {
	case <-ctx.Done():
		log.Println("ctx has been canceled")
	case <-time.After(2 * time.Second):
		assert.Fail(t, "context has not been canceled")
	}
}
