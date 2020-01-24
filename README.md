# gomulocity
Comulocity (c8y) go SDK.

## 'devicecontrol'-API

### Device credentials
Create newDeviceRequest:
```go
    newDeviceRequest, err := 'c8y>Client'.Devicecontrol.CreateNewDeviceRequest(<newDeviceRequestID>)
```
Find all newDeviceRequest (newDeviceRequestCollection)
```go
    newDeviceRequestCollection, err := 'c8y>Client'.Devicecontrol.NewDeviceRequestCollections(meta.Page(3))
```