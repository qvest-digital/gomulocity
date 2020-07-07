package generic

import (
	"encoding/json"
	"reflect"
	"testing"
)

type B struct {
	Foo    string                 `json:"foo"`
	Bar    int                    `json:"bar"`
	Baz    map[string]interface{} `jsonc:"flat"`
	FooBar string                 `jsonc:"flat"` // -> expect normal handling
}

type A struct {
	Bs []B    `json:"bList" jsonc:"collection"`
	C  int    `json:"c"`
	D  string `json:"d" jsonc:"collection"` // expect only `json:"d"` handling
}

const testJson = `{
		"bList":[
			{
				"foo":"Hallo",
				"bar":1,
				"FooBar": "myFooBar",
				"custom1":"#Custom1",
				"custom2":4711,
				"custom3": [
					"Hallo",
					"Welt"
				]
			},
			{
				"foo":"Hallo2",
				"bar":2,
				"FooBar": "myFooBar",
				"custom1":"#Custom1",
				"custom2":4711,
				"custom3": [
					"Hallo",		
					"Welt"
				]
			}
		],
		"c":4711,
		"d": "myDValue"
	}`

var additionalData = map[string]interface{}{
	"custom1": "#Custom1",
	"custom2": 4711,
	"custom3": []string{"Hallo", "Welt"},
}

var testObject = &A{C: 4711, D: "myDValue", Bs: []B{
	{Foo: "Hallo", Bar: 1, Baz: additionalData, FooBar: "myFooBar"},
	{Foo: "Hallo2", Bar: 2, Baz: additionalData, FooBar: "myFooBar"},
}}

func TestJsonc_Marshal_Lists(t *testing.T) {
	j, err := JsonFromObject(testObject)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	m := make(map[string]interface{})
	n := make(map[string]interface{})
	_ = json.Unmarshal([]byte(testJson), &m)
	_ = json.Unmarshal([]byte(j), &n)

	if !reflect.DeepEqual(n, m) {
		t.Errorf("JsonFromObject - json = %v, want %v", m, n)
	}
}

func TestJsonc_Unmarshal_Lists(t *testing.T) {
	a := &A{}
	err := ObjectFromJson([]byte(testJson), a)

	if err != nil {
		t.Errorf("ObjectFromJson - unexpected error %v", err)
	}

	if !reflect.DeepEqual(a, testObject) {
		t.Errorf("ObjectFromJson - json = %v \n want %v", a, testObject)
	}
}
