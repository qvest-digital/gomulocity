package managedObjects


func NewManagedObjectCollection_ResponseBody(next string) string {
	return `{
"next": "`+next+`",
"self": "https://t200588189.cumulocity.com/inventory/managedObjects?query=has(c8y_IsDevice)&pageSize=3&currentPage=1",
"managedObjects": [
{
"additionParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/additionParents",
"references": []
},
"owner": "device_4D8AFED3",
"childDevices": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/childDevices",
"references": []
},
"childAssets": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/childAssets",
"references": []
},
"creationTime": "2019-11-07T20:48:43.472Z",
"type": "c8y_SensorPhone",
"lastUpdated": "2020-04-24T09:36:05.112Z",
"childAdditions": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/childAdditions",
"references": [
{
"managedObject": {
"name": "Creates alarm when measurements are missing",
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/5755502",
"id": "5755502"
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/childAdditions/5755502"
}
]
},
"name": "Ridcully",
"assetParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/assetParents",
"references": []
},
"deviceParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353/deviceParents",
"references": []
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/2323353",
"id": "2323353",
"c8y_Firmware": {
"version": "13.2"
},
"c8y_Availability": {
"lastMessage": "2020-05-29T12:50:04.081Z",
"status": "UNAVAILABLE"
},
"com_cumulocity_model_Agent": {},
"c8y_ActiveAlarmsStatus": {
"major": 0,
"critical": 0
},
"c8y_IsDevice": [],
"c8y_RequiredAvailability": {
"responseInterval": 1
},
"c8y_Connection": {
"status": "DISCONNECTED"
},
"c8y_SupportedOperations": [
"c8y_Message",
"c8y_Relay"
],
"c8y_IsSensorPhone": [],
"c8y_Hardware": {
"serialNumber": "",
"model": "iPhone"
},
"c8y_DataPoint": {
"4955673840271746": {
"fragment": "c8y_Gyroscope",
"unit": "°/s",
"color": "#365381",
"series": "gyroY",
"lineType": "line",
"label": "c8y_Gyroscope => gyroY",
"_id": "4955673840271746",
"renderType": "min"
},
"17873621731249711": {
"fragment": "c8y_Gyroscope",
"unit": "°/s",
"color": "#365381",
"series": "gyroY",
"lineType": "line",
"label": "c8y_Gyroscope => gyroY",
"_id": "17873621731249711",
"renderType": "min"
},
"9807821862130062": {
"fragment": "c8y_Gyroscope",
"unit": "°/s",
"color": "#365381",
"series": "gyroY",
"lineType": "line",
"label": "c8y_Gyroscope => gyroY",
"_id": "9807821862130062",
"renderType": "min"
},
"7681935061014844": {
"fragment": "c8y_Gyroscope",
"unit": "°/s",
"color": "#365381",
"series": "gyroY",
"lineType": "line",
"label": "c8y_Gyroscope => gyroY",
"_id": "7681935061014844",
"renderType": "min"
}
}
},
{
"additionParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397/additionParents",
"references": []
},
"owner": "l.buhl@tarent.de",
"childDevices": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397/childDevices",
"references": []
},
"childAssets": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397/childAssets",
"references": []
},
"creationTime": "2019-08-27T04:37:28.768Z",
"lastUpdated": "2019-09-03T13:11:11.695Z",
"childAdditions": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397/childAdditions",
"references": []
},
"name": "flora2",
"assetParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397/assetParents",
"references": []
},
"deviceParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397/deviceParents",
"references": []
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232397",
"id": "232397",
"c8y_IsDevice": {},
"c8y_DataPoint": {
"1009058352217913": {
"fragment": "c8y_TemperatureMeasurement",
"unit": "°C",
"color": "#ba9d16",
"series": "temperature",
"lineType": "line",
"label": "Temperature",
"_id": "1009058352217913",
"renderType": "min"
}
}
},
{
"additionParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/additionParents",
"references": []
},
"owner": "l.buhl@tarent.de",
"childDevices": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childDevices",
"references": []
},
"childAssets": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAssets",
"references": []
},
"creationTime": "2019-08-27T04:37:24.074Z",
"lastUpdated": "2019-09-03T13:11:11.720Z",
"childAdditions": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAdditions",
"references": [
{
"managedObject": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/347746",
"id": "347746"
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAdditions/347746"
}
]
},
"name": "flora1",
"assetParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/assetParents",
"references": []
},
"deviceParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/deviceParents",
"references": []
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704",
"id": "232704",
"c8y_IsDevice": {},
"c8y_DataPoint": {
"9079903668235025": {
"fragment": "c8y_VoltageMeasurement",
"unit": "%",
"color": "#fe46d4",
"series": "battery",
"lineType": "line",
"label": "c8y_VoltageMeasurement => battery",
"_id": "9079903668235025",
"renderType": "min"
},
"61470130711723": {
"fragment": "c8y_TemperatureMeasurement",
"unit": "°C",
"color": "#bec591",
"series": "temperature",
"lineType": "line",
"label": "Temperature",
"_id": "61470130711723",
"renderType": "min"
},
"7425000382785754": {
"fragment": "c8y_TemperatureMeasurement",
"unit": "°C",
"color": "#3bc487",
"series": "temperature",
"lineType": "line",
"label": "c8y_TemperatureMeasurement => temperature",
"_id": "7425000382785754",
"renderType": "min"
}
}
}
],
"statistics": {
"currentPage": 1,
"pageSize": 3
}
}`
}
