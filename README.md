![Gomulocity](https://github.com/tarent/gomulocity/blob/master/docs/gomulocity.jpg?raw=true)

# Gomulocity REST SDK #

![Go](https://github.com/tarent/gomulocity/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tarent/gomulocity)](https://goreportcard.com/report/github.com/tarent/gomulocity)
![GitHub](https://img.shields.io/github/license/tarent/gomulocity)
![GitHub last commit](https://img.shields.io/github/last-commit/tarent/gomulocity)

 Gomulocity is a Go library to interact with the REST API of the Cumulocity IoT platform of Software AG (c8y).

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
## Table of Contents ##

- [Gomulocity REST SDK](#gomulocity-rest-sdk)
- [Usage](#usage)
    - [Device Bootstrap](#device-bootstrap)
        - [Configuration](#configuration)
        - [Device Registration API](#device-registration-api)
        - [Device Credentials API](#device-credentials-api)
- [Realtime Notification](#realtime-notification)
- [Feature coverage](#feature-coverage)
- [Contributing](#contributing)
- [License](#license)

<!-- markdown-toc end -->

# Usage #

APIs are split into seperate imports you can select from:

```go
import "github.com/tarent/gomulocity/alarm"
import "github.com/tarent/gomulocity/events"
import "github.com/tarent/gomulocity/measurement"
import "github.com/tarent/gomulocity/inventory"
```

The APIs need clients with credentials to work.

``` go
var c8yClient = &generic.Client{
	HTTPClient: http.DefaultClient,
	BaseURL:    "https://management.cumulocity.com",
	Username:   "user",
	Password:   "password",
}
```

## Device Bootstrap ##

### Configuration ###

The bootstrap API needs basic credentials to be able to register your client. Please contact your platform provider for the correct bootstrap credentials. Assure that the base url points to the correct platform instance.

``` go
var bootstrapClient = &generic.Client{
	HTTPClient: http.DefaultClient,
	BaseURL:    "https://management.cumulocity.com",
	Username:   "bootstrapuser",
	Password:   "password",
}
```

You can then register your device with a unique ID at your tenant:

``` go
	deviceId := "uniqueDeviceID"
	deviceCredentialsApi := device_bootstrap.NewDeviceCredentialsApi(bootstrapClient)

    deviceCredentials, _ := deviceCredentialsApi.Create(deviceId)
```

To register a device, you need to add the registration with the unique ID to your tenant via registration API or UI. More information about the registration cycle is available [in the device integration part of the c8y docs.](https://cumulocity.com/guides/device-sdk/rest/) After obtaining credentials for your device, you need to create the device itself as a managed object. Use the inventory API to accomplish this.

### Device Registration API ###
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

### Device Credentials API ###
Create DeviceCredentials:
```go
    deviceCredentials, err := gomulocity.DeviceCredentials.Create("123")
```
# Realtime Notification #
To use the Realtime-notification-API you need to import it with:

```go
import "github.com/tarent/gomulocity/realtimenotification"
```

To build and start the API use StartRealtimeNotificationsAPI() which requires a context.Context to work with as well as credentials and an adress to connect to.
credentials have to be provided in the following pattern: "tenantID/username:password"
```go
api, err := StartRealtimeNotificationsAPI(ctx,"mycumulocitytenant/myusername:mypassword", "myadress.cumulocity.com")
```

After this we can Subscribe and Unsubscribe to channels by using:
```go
api.DoSubscribe("operations/{deviceID})
api.DoUnsubscribe("measurement/{deviceID}")
```

To stop the API we can cancel the context, or use an OS.Interrupt Signal.

All answers from c8y are available in the ```api.ResponseFromPolling``` channel as raw json. You need to unmarshall it to the corresponding objects depending on your subscriptions.
# Feature coverage #

REST API:

- [x] inventory/managedObjects
- [x] measurement
- [x] alarm
- [x] event
- [x] deviceControl/operations
- [x] bootstrapping
- [x] identity
- [x] Realtime notifications via websockets
- [x] audit
- [x] user
- [ ] tenant

# Contributing #

When contributing to this repository, please first discuss the change you wish to make via issue with the owners of this repository before making a change.

# License #

See [LICENSE file](LICENSE "LICENSE file").
