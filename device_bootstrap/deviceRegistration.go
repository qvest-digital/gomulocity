package device_bootstrap

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

const (
	DEVICE_REGISTRATION_API_PATH        = "/devicecontrol/newDeviceRequests"
	DEVICE_REGISTRATION_TYPE            = "application/vnd.com.nsn.cumulocity.NewDeviceRequest+json"
	DEVICE_REGISTRATION_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.newDeviceRequestCollection+json"
)

type Status string

const (
	WAITING_FOR_CONNECTION Status = "WAITING_FOR_CONNECTION"
	PENDING_ACCEPTANCE     Status = "PENDING_ACCEPTANCE"
	ACCEPTED               Status = "ACCEPTED"
)

/*
DeviceRegistration represent cumulocity's 'application/vnd.com.nsn.cumulocity.NewDeviceRequest+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#newdevicerequest-application-vnd-com-nsn-cumulocity-newdevicerequest-json
*/
type DeviceRegistration struct {
	Id               string      `json:"id,omitempty"`
	Status           Status      `json:"status,omitempty"`
	Self             string      `json:"self,omitempty"`
	Owner            string      `json:"owner,omitempty"`
	CustomProperties interface{} `json:"customProperties,omitempty"`
	CreationTime     *time.Time   `json:"creationTime,omitempty"`
	TenantId         string      `json:"tenantId,omitempty"`
}

/*
DeviceRequestCollection represent cumulocity's 'application/vnd.com.nsn.cumulocity.newDeviceRequestCollection+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#newdevicerequestcollection-application-vnd-com-nsn-cumulocity-newdevicerequestcollection-json
*/
type DeviceRegistrationCollection struct {
	Self                string                    `json:"self"`
	DeviceRegistrations []DeviceRegistration      `json:"newDeviceRequests"`
	Statistics          *generic.PagingStatistics `json:"statistics"`
	Prev                string                    `json:"prev"`
	Next                string                    `json:"next"`
}
