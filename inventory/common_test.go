package inventory

import (
	"github.com/tarent/gomulocity/generic"
	"net/http/httptest"
	"time"
)

const (
	USER     = "foo"
	PASSWORD = "bar"
)

var creationTime, _ = time.Parse(time.RFC3339, "2020-07-03T10:16:35.870+02:00")

func buildInventoryApi(testServer *httptest.Server) InventoryApi {
	c := &generic.Client{
		HTTPClient: testServer.Client(),
		BaseURL:    testServer.URL,
		Username:   USER,
		Password:   PASSWORD,
	}

	return NewInventoryApi(c)
}

func buildInventoryReferenceApi(testServer *httptest.Server) InventoryReferenceApi {
	c := &generic.Client{
		HTTPClient: testServer.Client(),
		BaseURL:    testServer.URL,
		Username:   USER,
		Password:   PASSWORD,
	}

	return NewInventoryReferenceApi(c)
}

var newManagedObject = &NewManagedObject{
	Type:         "test-type",
	Name:         "Test Device",
	CreationTime: &creationTime,
}

var managedObjectUpdate = &ManagedObjectUpdate{
	Type: "updated test-type",
	Name: "updated Test Device",
}
var expectedUpdateRequestBody = `{"type":"updated test-type","name":"updated Test Device"}`

var expectedRequestBody = `{"type":"test-type","name":"Test Device","creationTime":"2020-07-03T10:16:35.87+02:00"}`
var managedObjectId = "9963944"
var referenceId = "4711"
var query = "$filter=name eq '*Test*' $orderby=id desc"

var inventoryFilter = &InventoryFilter{
	Type:         "test-type",
	FragmentType: "Test-FragmentType",
	Ids:          []string{managedObjectId},
	Text:         "Test Device",
}
var expectedQuery = "fragmentType=Test-FragmentType&ids=9963944&pageSize=5&text=Test+Device&type=test-type"

var expectedManagedObject = &ManagedObject{
	Id:              "9963944",
	Type:            "test-type",
	Name:            "Test Device",
	CreationTime:    creationTime,
	LastUpdated:     creationTime,
	Self:            "https://t0815.cumulocity.com/inventory/managedObjects/9963944",
	Owner:           "gomulocity",
	AdditionParents: AdditionParents{References: []interface{}{}, Self: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/additionParents"},
	AssetParents:    AssetParents{References: []interface{}{}, Self: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/assetParents"},
	DeviceParents:   DeviceParents{References: []interface{}{}, Self: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/deviceParents"},
	ChildAdditions: ChildAdditions{References: []struct {
		ManagedObject struct {
			Id   string "json:\"id,omitempty\""
			Name string "json:\"name,omitempty\""
			Self string "json:\"self,omitempty\""
		} "json:\"managedObject,omitempty\""
		Self string "json:\"self,omitempty\""
	}{}, Self: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAdditions"},
	ChildAssets:  ChildAssets{References: []interface{}{}, Self: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAssets"},
	ChildDevices: ChildDevices{References: []interface{}{}, Self: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices"},
	AdditionalFields: map[string]interface{}{
		"custom": "hello",
	},
}

var givenResponseBody = `{
				"id": "9963944",
				"type": "test-type",
				"name": "Test Device",
				"creationTime": "2020-07-03T10:16:35.870+02:00",
				"lastUpdated": "2020-07-03T10:16:35.870+02:00",
				"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944",
				"owner": "gomulocity",
				"additionParents": {
					"references": [],
					"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/additionParents"
				},
				"assetParents": {
					"references": [],
					"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/assetParents"
				},
				"deviceParents": {
					"references": [],
					"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/deviceParents"
				},
				"childAdditions": {
					"references": [],
					"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAdditions"
				},
				"childAssets": {
					"references": [],
					"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAssets"
				},
				"childDevices": {
					"references": [],
					"self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices"
				},
				"custom":"hello"
			}`

var expectedManagedObjectCollection = &ManagedObjectCollection{
	Self:           "https://t0815.cumulocity.com/inventory/managedObjects?ids=9963944&text=Test%20Device&type=integration-test&pageSize=5&currentPage=1",
	ManagedObjects: []ManagedObject{*expectedManagedObject},
	Statistics: &generic.PagingStatistics{
		PageSize:    5,
		CurrentPage: 1,
	},
	Next: "https://t0815.cumulocity.com/inventory/managedObjects?ids=9963944&text=Test%20Device&type=integration-test&pageSize=5&currentPage=2",
}

var givenManagedObjectCollectionResponse = `{
    "self": "https://t0815.cumulocity.com/inventory/managedObjects?ids=9963944&text=Test%20Device&type=integration-test&pageSize=5&currentPage=1",
    "managedObjects": [
        {
            "id": "9963944",
            "type": "test-type",
            "name": "Test Device",
            "creationTime": "2020-07-03T10:16:35.870+02:00",
            "lastUpdated": "2020-07-03T10:16:35.870+02:00",
            "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944",
            "owner": "gomulocity",
            "additionParents": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/additionParents"
            },
            "assetParents": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/assetParents"
            },
            "deviceParents": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/deviceParents"
            },
            "childAdditions": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAdditions"
            },
            "childAssets": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAssets"
            },
            "childDevices": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices"
            },
			"custom":"hello"
        }
    ],
    "statistics": {
        "pageSize": 5,
        "currentPage": 1
    },
    "next": "https://t0815.cumulocity.com/inventory/managedObjects?ids=9963944&text=Test%20Device&type=integration-test&pageSize=5&currentPage=2"
}`

var managedObjectCollectionTemplate = `{
    "self": "https://t0815.cumulocity.com/inventory/managedObjects?type=test-type&pageSize=5&currentPage=1",
	"next": "https://t0815.cumulocity.com/inventory/managedObjects?type=test-type&pageSize=5&currentPage=2",
	"managedObjects": [%s],
    "statistics": {
			"pageSize": 5,
			"currentPage": 1
	}
}`

var childManagedObject = &ManagedObject{
	Id:   "4711",
	Self: "https://t0815.cumulocity.com/inventory/managedObjects/4711",
}

var expectedManagedObjectReference = &ManagedObjectReference{
	ManagedObject: *childManagedObject,
	Self:          "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices/4711",
}

var expectedNewReferenceRequestBody = `{"managedObject":{"id":"4711"}}`

var givenReferenceResponseBody = `{"managedObject":{"id":"4711","self":"https://t0815.cumulocity.com/inventory/managedObjects/4711"},"self":"https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices/4711"}`

var givenManagedObjectReferenceCollectionResponse = `{
    "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices?pageSize=5&currentPage=1",
    "references": [
        {
            "managedObject": {
                "id": "4711",
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/4711"
            },
            "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices/4711"
        }
    ],
    "statistics": {
        "pageSize": 5,
        "currentPage": 1
    },
    "next": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices?pageSize=5&currentPage=2"
}`

var expectedManagedObjectReferenceCollection = &ManagedObjectReferenceCollection{
	Self:       "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices?pageSize=5&currentPage=1",
	References: []ManagedObjectReference{*expectedManagedObjectReference},
	Statistics: &generic.PagingStatistics{
		PageSize:    5,
		CurrentPage: 1,
	},
	Prev: "",
	Next: "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices?pageSize=5&currentPage=2",
}

var managedObjectReferenceCollectionTemplate = `{
    "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices?&pageSize=5&currentPage=1",
	"next": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices?&pageSize=5&currentPage=2",
	"references": [%s],
    "statistics": {
			"pageSize": 5,
			"currentPage": 1
	}
}`

var givenManagedObjectCollectionResponseWithAdditionalFields = `{
    "self": "https://t0815.cumulocity.com/inventory/managedObjects?ids=9963944&text=Test%20Device&type=integration-test&pageSize=5&currentPage=1",
    "managedObjects": [
        {
            "id": "9963944",
            "type": "test-type",
            "name": "Test Device",
            "creationTime": "2020-07-03T10:16:35.870+02:00",
            "lastUpdated": "2020-07-03T10:16:35.870+02:00",
            "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944",
            "owner": "gomulocity",
            "additionParents": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/additionParents"
            },
            "assetParents": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/assetParents"
            },
            "deviceParents": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/deviceParents"
            },
            "childAdditions": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAdditions"
            },
            "childAssets": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childAssets"
            },
            "childDevices": {
                "references": [],
                "self": "https://t0815.cumulocity.com/inventory/managedObjects/9963944/childDevices"
            },
			"custom1":"hello",
			"custom2":{
				"foo":"bar"
			}
        }
    ],
    "statistics": {
        "pageSize": 5,
        "currentPage": 1
    },
    "next": "https://t0815.cumulocity.com/inventory/managedObjects?ids=9963944&text=Test%20Device&type=integration-test&pageSize=5&currentPage=2"
}`
