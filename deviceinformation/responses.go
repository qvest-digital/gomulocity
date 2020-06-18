package deviceinformation

const responseBodyDeviceInformation = `{
"next": "https://t200588189.cumulocity.com/inventory/managedObjects?pageSize=1&currentPage=2",
"self": "https://t200588189.cumulocity.com/inventory/managedObjects?pageSize=1&currentPage=1",
"managedObjects": [
{
"additionParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/additionParents",
"references": []
},
"childDevices": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/childDevices",
"references": []
},
"childAssets": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/childAssets",
"references": []
},
"type": "c8y_Application_2835",
"childAdditions": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/childAdditions",
"references": []
},
"name": "device-simulator",
"assetParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/assetParents",
"references": []
},
"deviceParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/deviceParents",
"references": []
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100",
"id": "100",
"c8y_Status": {
"lastUpdated": {
"date": {
"$date": "2020-06-03T17:11:59.003Z"
},
"offset": 0
},
"instances": {
"device-simulator-scope-management-deployment-7d56bb749c-wg64m": {
"lastUpdated": {
"date": {
"$date": "2020-06-03T17:11:59.003Z"
},
"offset": 0
},
"memoryInBytes": 1073741824,
"restarts": 1,
"cpuInMillis": 2000
}
},
"details": {
"desired": 1,
"aggregatedResources": {
"memory": "1073M",
"cpu": "2000m"
},
"active": 1,
"restarts": 1
},
"status": "Up"
},
"applicationOwner": "management",
"applicationId": "2835"
}
],
"statistics": {
"totalPages": 765,
"currentPage": 1,
"pageSize": 1
}
}`