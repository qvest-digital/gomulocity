# gomulocity
Cumulocity (c8y) go SDK.

# Usage
Example:
```go
import (
	"github.com/tarent/gomulocity"
)

func main() {
	gomulocity := gomulocity.NewGomulocity("https://<tenant>.<c8yHost>", "<username>", "<password>", "<bootstrap-user>", "<bootstrap-password>")
}
```

## Device Bootstrap

### Device Registration API
Start a new device registration with a unique device ID:
```go
    deviceRegistration, err := gomulocity.DeviceRegistration.Create("123")
```
Get a device registration by device ID:
```go
    deviceRegistration, err := gomulocity.DeviceRegistration.Get("123")
```
Get all device registrations page by page:
```go
    deviceRegistrations, err := gomulocity.DeviceRegistration.GetAll(10)
    deviceRegistrations, err = gomulocity.DeviceRegistration.NextPage(deviceRegistrations)
    deviceRegistrations, err = gomulocity.DeviceRegistration.PreviousPage(deviceRegistrations)
```
Update device registration status:
```go
    deviceRegistration, err := gomulocity.DeviceRegistration.Update("123", device_bootstrap.ACCEPTED)
```
Delete device registration by device ID:
```go
    err := gomulocity.DeviceRegistration.Delete("123")
```

### Device Credentials API
Create DeviceCredentials:
```go
    deviceCredentials, err := gomulocity.DeviceCredentials.Create("123")
```
