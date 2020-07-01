package device_bootstrap

const (
	DEVICE_CREDENTIALS_API_PATH = "/devicecontrol/deviceCredentials"
	DEVICE_CREDENTIALS_TYPE     = "application/vnd.com.nsn.cumulocity.deviceCredentials+json"
)

/*
DeviceCredentials represent cumulocity's 'application/vnd.com.nsn.cumulocity.deviceCredentials+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#devicecredentials-application-vnd-com-nsn-cumulocity-devicecredentials-json
*/
type DeviceCredentials struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Self     string `json:"self,omitempty"`
}
