# gomulocity
Comulocity (c8y) go SDK.

# Usage
Example:
```go
import (
	"fmt"
	"github.com/tarent/gomulocity/pkg/c8y/devicecontrol"
	"github.com/tarent/gomulocity/pkg/c8y/meta"
	"net/http"
	"testing"
)

func main() {
	c := devicecontrol.Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://<tenant>.<c8yHost>",
		Username:   "<username>",
		Password:   "<password>",
	}
    
    c.Devicecontrol //see: 'devicecontrol'-API 
}
```

## 'devicecontrol'-API

### Device credentials
Create newDeviceRequest:
```go
    newDeviceRequest, err := 'c8y.Client'.Devicecontrol.CreateNewDeviceRequest(<newDeviceRequestID>)
```
Find all newDeviceRequest (newDeviceRequestCollection):
```go
    newDeviceRequestCollection, err := 'c8y.Client'.Devicecontrol.NewDeviceRequestCollections(meta.Page(3))
```
Update newDeviceRequest:
```go
    err := 'c8y.Client'.Devicecontrol.UpdateNewDeviceRequest(<newDeviceRequestID>, <newDeviceRequestStatus>)
```
Delete newDeviceRequest:
```go
    err := 'c8y.Client'.Devicecontrol.DeleteNewDeviceRequest(<newDeviceRequestID>)
```