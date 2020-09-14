package realtimenotification

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_RealtimeSubscription_Handshake_send(t *testing.T) {
	send := make(chan []byte, 5)
	defer close(send)
	login := Login{
		Authentification: Auth{
			Token: "somecredentials",
		},
		SystemOfUnits: "metric",
	}
	api := RealtimeNotificationAPI{send: send, login: login}

	api.DoHandshake()

	request := <-api.send

	require.NotEmpty(t, request)
	assert.Equal(t, `{"channel":"/meta/handshake","ext":{"com.cumulocity.authn":{"token":"somecredentials"},"systemOfunits":"metric"},"version":"1.0","minimumVersion":"1.0","supportedConnectionTypes":["websocket"],"advice":{"timeout":60000,"interval":10000}}`, string(request))

}
