package inventory

const ReferenceCollectionJson = `{
    "next": "https://t200588189.cumulocity.com/inventory/managedObjects/104940/childAdditions?pageSize=5&currentPage=2",
    "references": [
        {
            "managedObject": {
                "additionParents": {
                    "references": [],
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167/additionParents"
                },
                "assetParents": {
                    "references": [],
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167/assetParents"
                },
                "c8y_Dashboard": {
                    "children": {
                        "5879719650867943": {
                            "_height": 6,
                            "_width": 10,
                            "_x": 0,
                            "_y": 0,
                            "classes": {
                                "panel-content-dark": true,
                                "panel-title-regular": true
                            },
                            "col": 6,
                            "config": {
                                "alarmsEventsConfigs": [],
                                "datapoints": [{}],
                                "dateFrom": "2020-03-02T14:34:10.000Z",
                                "dateTo": "2020-03-02T15:34:10.000Z",
                                "interval": "hours",
                                "realtime": true
                            },
                            "configTemplateUrl": "dataPointExplorer/views/widgetConfig.html",
                            "id": "5879719650867943",
                            "name": "Data points graph",
                            "position": 0,
                            "templateUrl": "dataPointExplorer/views/widget.html",
                            "title": "Data points graph"
                        }
                    },
                    "classes": {
                        "dashboard-theme-dark": true
                    },
                    "global": false,
                    "icon": "th",
                    "isFrozen": false,
                    "name": "Dashboard",
                    "priority": 10000,
                    "translateWidgetTitle": false,
                    "widgetClasses": {
                        "panel-title-regular": true
                    }
                },
                "c8y_Dashboard!group!104940": {},
                "childAdditions": {
                    "references": [],
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167/childAdditions"
                },
                "childAssets": {
                    "references": [],
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167/childAssets"
                },
                "childDevices": {
                    "references": [],
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167/childDevices"
                },
                "creationTime": "2019-08-27T05:56:42.484Z",
                "deviceParents": {
                    "references": [],
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167/deviceParents"
                },
                "id": "235167",
                "lastUpdated": "2020-03-02T15:34:51.747Z",
                "owner": "l.buhl@tarent.de",
                "self": "https://t200588189.cumulocity.com/inventory/managedObjects/235167"
            },
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/104940/childAdditions/235167"
        }
    ],
    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/104940/childAdditions?pageSize=5&currentPage=1",
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`

const ReferenceByID = `{
    "managedObject": {
        "additionParents": {
            "references": [],
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/additionParents"
        },
        "assetParents": {
            "references": [],
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/assetParents"
        },
        "c8y_DataPoint": {
            "61470130711723": {
                "_id": "61470130711723",
                "color": "#bec591",
                "fragment": "c8y_TemperatureMeasurement",
                "label": "Temperature",
                "lineType": "line",
                "renderType": "min",
                "series": "temperature",
                "unit": "°C"
            },
            "7425000382785754": {
                "_id": "7425000382785754",
                "color": "#3bc487",
                "fragment": "c8y_TemperatureMeasurement",
                "label": "c8y_TemperatureMeasurement => temperature",
                "lineType": "line",
                "renderType": "min",
                "series": "temperature",
                "unit": "°C"
            },
            "9079903668235025": {
                "_id": "9079903668235025",
                "color": "#fe46d4",
                "fragment": "c8y_VoltageMeasurement",
                "label": "c8y_VoltageMeasurement => battery",
                "lineType": "line",
                "renderType": "min",
                "series": "battery",
                "unit": "%"
            }
        },
        "c8y_IsDevice": {},
        "childAdditions": {
            "references": [
                {
                    "managedObject": {
                        "id": "347746",
                        "self": "https://t200588189.cumulocity.com/inventory/managedObjects/347746"
                    },
                    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAdditions/347746"
                }
            ],
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAdditions"
        },
        "childAssets": {
            "references": [],
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAssets"
        },
        "childDevices": {
            "references": [],
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childDevices"
        },
        "creationTime": "2019-08-27T04:37:24.074Z",
        "deviceParents": {
            "references": [],
            "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/deviceParents"
        },
        "id": "232704",
        "lastUpdated": "2019-09-03T13:11:11.720Z",
        "name": "flora1",
        "owner": "l.buhl@tarent.de",
        "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704"
    },
    "self": "https://t200588189.cumulocity.com/inventory/managedObjects/232704/childAdditions/232704"
}`

const ManagedObjectByID = `{
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
"deviceParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/deviceParents",
"references": []
},
"assetParents": {
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100/assetParents",
"references": []
},
"self": "https://t200588189.cumulocity.com/inventory/managedObjects/100",
"id": "100",
"c8y_Status": {
"lastUpdated": {
"date": {
"$date": "2020-07-06T11:54:01.015Z"
},
"offset": 0
},
"instances": {
"device-simulator-scope-management-deployment-77678578b4-vkn66": {
"lastUpdated": {
"date": {
"$date": "2020-07-06T11:54:01.015Z"
},
"offset": 0
},
"memoryInBytes": 1073741824,
"restarts": 2,
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
"restarts": 2
},
"status": "Up"
},
"applicationOwner": "management",
"applicationId": "2835"
}`

var NewManagedObjectCollection_ResponseBody = func(next string) string {
	return `{
"next": "` + next + `",
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

var UpdatedManagedObject = func(name string) string {
	return `{
    "id" : "104940",
    "name" : "` + name + `",
    "self" : "https://t200588189.cumulocity.com/inventory/managedObjects/104940",
    "type" :"c8y_DeviceGroup",
    "lastUpdated": "2019-08-23T15:10:00.653Z",
    "com_othercompany_StrongTypedClass" : {},
    "childDevices": {
		"self":"https://t200588189.cumulocity.com/inventory/managedObjects/104940/childDevices",
		"references":[{}]
	}
  }`
}
