# gomulocity
Comulocity (c8y) go SDK.

# Usage
Example:
```go
import (
	"github.com/tarent/gomulocity"
)

func main() {
	client := gomulocity.NewClient("https://<tenant>.<c8yHost>", "<username>", "<password>", "<bootstrap-user>", "<bootstrap-password>")
}
```

## Device Bootstrap

### Device Registration API
Start a new device registration with a unique device ID:
```go
    deviceRegistration, err := client.DeviceRegistration.Create("123")
```
Get a device registration by device ID:
```go
    deviceRegistration, err := deviceRegistrationApi.Get("123")
```
Get all device registrations page by page:
```go
    deviceRegistrations, err := deviceRegistrationApi.GetAll(10)
    deviceRegistrations, err = deviceRegistrationApi.NextPage(deviceRegistrations)
    deviceRegistrations, err = deviceRegistrationApi.PreviousPage(deviceRegistrations)
```
Update device registration status:
```go
    deviceRegistration, err := deviceRegistrationApi.Update("123", device_bootstrap.ACCEPTED)
```
Delete device registration by device ID:
```go
    err := deviceRegistrationApi.Delete("123")
```

### Device Credentials API
Create DeviceCredentials:
```go
    deviceCredentials, err := deviceCredentialsApi.Create("123")
```