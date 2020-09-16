package realtimenotification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_RealtimeSubscription_Handshake_send_and_receive(t *testing.T) {
	send := make(chan []byte, 5)
	defer close(send)
	response := make(chan []byte, 5)
	defer close(response)
	login := Login{
		Authentification: Auth{
			Token: "somecredentials",
		},
		SystemOfUnits: "metric",
	}

	api := RealtimeNotificationAPI{
		send:     send,
		response: response,
		login:    login,
		timeout:  5 * time.Second,
	}
	answer := []AuthResponse{
		AuthResponse{
			Id:         "someId",
			ClientId:   "someClientId",
			Successful: true,
		},
	}

	jsonresponse, err := json.Marshal(answer)
	require.NoError(t, err)
	api.response <- jsonresponse

	err = api.doHandshake()
	require.NoError(t, err)

	request := <-api.send

	require.NotEmpty(t, request)
	assert.Equal(t, `{"channel":"/meta/handshake","ext":{"com.cumulocity.authn":{"token":"somecredentials"},"systemOfunits":"metric"},"version":"1.0","minimumVersion":"1.0","supportedConnectionTypes":["websocket"],"advice":{"timeout":60000,"interval":10000}}`, string(request))

}

func Test_RealtimeSubscription_Handshake_times_out(t *testing.T) {
	send := make(chan []byte, 5)
	defer close(send)
	login := Login{
		Authentification: Auth{
			Token: "somecredentials",
		},
		SystemOfUnits: "metric",
	}

	api := RealtimeNotificationAPI{
		timeout: 0 * time.Second,
		send:    send,
		login:   login,
	}

	err := api.doHandshake()
	require.Error(t, err)
}
