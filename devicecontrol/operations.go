package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/url"
	"time"
)

const (
	OPERATION_CONTENT_TYPE       = "application/vnd.com.nsn.cumulocity.operation+json"
	OPERATION_ACCEPT_HEADER      = "application/vnd.com.nsn.cumulocity.operation+json"
	BULK_OPERATION_CONTENT_TYPE  = "application/vnd.com.nsn.cumulocity.bulkOperation+json"
	BULK_OPERATION_ACCEPT_HEADER = "application/vnd.com.nsn.cumulocity.bulkOperation+json"
)

type OperationCollection struct {
	Self       string                    `json:"self"`
	Operations []Operation               `json:"operations"`
	Statistics *generic.PagingStatistics `json:"statistics"`
	Prev       string                    `json:"prev"`
	Next       string                    `json:"next"`
}

type Operation struct {
	Self          string    `json:"self"`
	DeviceID      string    `json:"deviceId"`
	DeviceName    string    `json:"deviceName"`
	OperationID   string    `json:"id"`
	CreationTime  time.Time `json:"creationTime"`
	Status        string    `json:"status"`
	FailureReason string    `json:"failureReason"`
	Description   string    `json:"description"`
	Delivery      struct {
		Time   time.Time `json:"time"`
		Status string    `json:"status"`
		Log    []struct {
			Time   time.Time `json:"time"`
			Status string    `json:"status"`
		} `json:"log"`
	} `json:"delivery"`
	AdditionalFields map[string]interface{} `jsonc:"flat"` //c8y_Command/c8y_Restart, ...
}

/*
"delivery": {
        "time": "2020-05-14T10:18:33.880Z",
        "status": "DELIVERED",
        "log": [
            {
                "time": "2020-05-14T10:18:33.649Z",
                "status": "PENDING"
            },
            {
                "time": "2020-05-14T10:18:33.679Z",
                "status": "SEND"
            }
        ]
    },
*/

type NewOperation struct {
	DeviceID         string                 `json:"deviceId"`
	AdditionalFields map[string]interface{} `jsonc:"flat"`
}

type UpdateOperation struct {
	Status           string                 `json:"status"`
	AdditionalFields map[string]interface{} `jsonc:"flat"`
}

type NewBulkOperation struct {
	StartDate          time.Time              `json:"startDate"`
	CreationRamp       int                    `json:"creationRamp"`
	OperationPrototype map[string]interface{} `jsonc:"flat"`
}

type BulkOperationCollection struct {
	Next           string                   `json:"next"`
	Self           string                   `json:"self"`
	Prev           string                   `json:"prev"`
	BulkOperations []BulkOperation          `json:"bulkOperations"`
	Statistics     generic.PagingStatistics `json:"statistics"`
}

type UpdateBulkOperation struct {
	CreationRamp int `json:"creationRamp"`
}

type BulkOperation struct {
	CreationRamp int    `json:"creationRamp"`
	GroupID      string `json:"groupId"`
	Description  string `json:"description"`
	Progress     struct {
		All        int `json:"all"`
		Executing  int `json:"executing"`
		Failed     int `json:"failed"`
		Pending    int `json:"pending"`
		Successful int `json:"successful"`
	} `json:"progress"`
	Self               string                 `json:"self"`
	ID                 int                    `json:"id"`
	StartDate          time.Time              `json:"startDate"`
	Status             string                 `json:"status"`
	OperationPrototype map[string]interface{} `jsonc:"flat"`
}

type OperationQuery struct {
	DeviceID          string
	Status            string
	AgentID           string
	AgentIDAndStatus  string
	DeviceIDAndStatus string
}

func (o *OperationQuery) QueryParams(params *url.Values) {
	if len(o.DeviceIDAndStatus) > 0 {
		params.Add("operationsByDeviceIdAndStatus", o.DeviceIDAndStatus)
	}

	if len(o.AgentIDAndStatus) > 0 {
		params.Add("operationsByAgentIdAndStatus", o.AgentIDAndStatus)
	}

	if len(o.AgentIDAndStatus) == 0 && len(o.DeviceIDAndStatus) == 0 {
		if len(o.DeviceID) > 0 {
			params.Add("operationsByDeviceId", o.DeviceID)
		}

		if len(o.Status) > 0 {
			params.Add("operationsByStatus", o.Status)
		}

		if len(o.AgentID) > 0 {
			params.Add("operationsByAgentId", o.AgentID)
		}
	}
}
