package devicecontrol

import (
	"encoding/json"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var creationTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")
var responseOperation = &Operation{
	DeviceID:         "123",
	DeviceName:       "name",
	OperationID:      "1",
	CreationTime:     creationTime,
	Status:           "status",
	FailureReason:    "",
	Description:      "description",
	Self:             "https://t200588189.cumulocity.com/devicecontrol/operations/4788662",
	AdditionalFields: map[string]interface{}{},
}

func buildOperationApi(url string) DeviceControl {
	httpClient := http.DefaultClient
	client := generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewDeviceControlApi(&client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func updateOperationHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var operation UpdateOperation
		_ = json.Unmarshal(body, &operation)
		updateOperationCapture = &operation
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal(responseOperation)
		_, _ = w.Write(response)
	}))
}

func createOperationHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

var requestCapture *http.Request
var createOperationCapture *NewOperation
var updateOperationCapture *UpdateOperation
var updateUrlCapture string

var operationID = "1111111"
var bulkOperationID = "1"
var operation = `{
    "status": "SUCCESSFUL",
    "self": "https://t200588189.cumulocity.com/devicecontrol/operations/4788662",
    "id": "1111111",
    "deviceId": "4788195",
    "description": "Restart device",
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
    "creationTime": "2020-05-14T10:18:33.477Z",
    "c8y_Restart": {},
	"custom1": "Hello"
}`
var operationCollectionTemplate = `{
    "next": "https://t0818.cumulocity.com/devicecontrol/operations/1111111?pageSize=5&currentPage=2",
    "self": "https://t0815.cumulocity.com/devicecontrol/operations/1111111?pageSize=5&currentPage=1",
    "operations": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`

var erroneousResponse = `{
    "error": "devicecontrol/Not Found",
    "message": "Finding device data from database failed : No operation for gid '47886623434324'!",
    "info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}`

var erroneousResponseBulkOperation = `{
    "error": "devicecontrol/Not Found",
    "info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
    "message": "Finding bulk operation from database failed : Cannot find document with ID ID [type=com_cumulocity_model_idtype_GId, value=4]!"
}`

var newOperation = `{
    "custom": "hello",
    "deviceId": "4788195"
}`

var bulkOperation = `{
    "creationRamp": 15,
    "groupId": "2323098",
    "id": 1,
    "operationPrototype": {
        "c8y_Command": {
            "text": "test"
        },
        "deliveryType": "SMS",
        "description": "Execute shell command"
    },
    "progress": {
        "all": 1,
        "executing": 0,
        "failed": 0,
        "pending": 0,
        "successful": 0
    },
    "self": "http://t200588189.cumulocity.com/devicecontrol/bulkoperations/1",
    "startDate": "2020-01-23T12:29:35.387Z",
    "status": "COMPLETED"
}`

var operationCollection = `{
    "self": "https://t200588189.cumulocity.com/devicecontrol/operations?pageSize=5&currentPage=1",
    "next": "https://t200588189.cumulocity.com/devicecontrol/operations?pageSize=5&currentPage=2",
    "statistics": {
        "pageSize": 5,
        "totalPages": 1,
        "currentPage": 1
    },
    "operations": [
        {
            "id": "3563887",
            "self": "https://t200588189.cumulocity.com/devicecontrol/operations/3563887",
            "status": "FAILED",
            "deviceId": "3249117",
            "deviceName": "Gate 154024",
            "description": "Command description",
            "c8y_Command": {
                "text": "<command>"
            },
            "creationTime": "2020-04-01T18:14:20.204Z",
            "failureReason": "Operation cancelled by user."
        },
        {
            "id": "3576460",
            "self": "https://t200588189.cumulocity.com/devicecontrol/operations/3576460",
            "status": "PENDING",
            "deviceId": "3249117",
            "deviceName": "Gate 154024",
            "description": "Command description",
            "c8y_Command": {
                "text": "<command>"
            },
            "creationTime": "2020-04-03T09:15:24.061Z"
        },
        {
            "id": "4788662",
            "self": "https://t200588189.cumulocity.com/devicecontrol/operations/4788662",
            "status": "SUCCESSFUL",
            "deviceId": "4788195",
            "delivery": {
                "log": [
                    {
                        "time": "2020-05-14T10:18:33.649Z",
                        "status": "PENDING"
                    },
                    {
                        "time": "2020-05-14T10:18:33.679Z",
                        "status": "SEND"
                    }
                ],
                "time": "2020-05-14T10:18:33.880Z",
                "status": "DELIVERED"
            },
            "deviceName": "My Java MQTT device",
            "description": "Restart device",
            "c8y_Restart": {},
            "creationTime": "2020-05-14T10:18:33.477Z"
        }
    ]
}`

var bulkOperationCollection = `{
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    },
    "next": "https://t200588189.cumulocity.com/devicecontrol/bulkoperations/?pageSize=5&currentPage=2",
    "self": "https://t200588189.cumulocity.com/devicecontrol/bulkoperations/?pageSize=5&currentPage=1",
    "bulkOperations": [
        {
            "progress": {
                "executing": 0,
                "failed": 0,
                "pending": 0,
                "successful": 0,
                "all": 1
            },
            "operationPrototype": {
                "deliveryType": "SMS",
                "description": "Execute shell command",
                "c8y_Command": {
                    "text": "test"
                }
            },
            "id": 1,
            "creationRamp": 15,
            "self": "https://t200588189.cumulocity.com/devicecontrol/bulkoperations/1",
            "status": "COMPLETED",
            "groupId": "2323098",
            "startDate": "2020-01-23T12:29:35.387Z"
        },
        {
            "progress": {
                "executing": 0,
                "failed": 0,
                "pending": 0,
                "successful": 2,
                "all": 2
            },
            "operationPrototype": {
                "deliveryType": "SMS",
                "description": "Execute shell command",
                "c8y_Command": {
                    "text": "test"
                }
            },
            "id": 12,
            "creationRamp": 14,
            "self": "https://t200588189.cumulocity.com/devicecontrol/bulkoperations/1",
            "status": "COMPLETED",
            "groupId": "2323098",
            "startDate": "2020-01-21T12:29:35.387Z"
        }
    ]
}`
