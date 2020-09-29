package realtimenotification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//Test_RealtimeNotification_simple needs working Credentials
func Test_RealtimeNotification_Simple(t *testing.T) {
	t.Log("test started")
	ctx := context.Background()

	//pattern: tennant/user:password
	credentials := ""
	if credentials ==""{
		t.Log("no credentials, aborting without failing as this is just a Demo")
		return
	}
	api, err := StartRealtimeNotificationsAPI(ctx, credentials, "https://management.cumulocity.com")
	require.NoError(t, err)
	t.Log("startup successful")

	err = api.DoSubscribe("operations/2867")
	require.NoError(t, err)
	t.Log("subscription successful")

	for {
		select {
		case answer := <-api.ResponseFromPolling:
			t.Log(string(answer))
		case <-time.After(10 * time.Second):
			api.stopPolling()
			return
		}
	}

}
