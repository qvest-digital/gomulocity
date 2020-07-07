package generic

import (
	"encoding/json"
	"reflect"
	"testing"
)

type B struct {
	Foo string                 `json:"foo"`
	Bar int                    `json:"bar"`
	Baz map[string]interface{} `jsonc:"flat"`
}

type A struct {
	Bs []B `json:"bList" jsonc:"collection"`
	C  int `json:"c"`
}

const testJson = `{
		"bList":[
			{
				"foo":"Hallo",
				"bar":1,
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
				"custom1":"#Custom1",
				"custom2":4711,
				"custom3": [
					"Hallo",		
					"Welt"
				]
			}
		],
		"c":4711
	}`

var additionalData = map[string]interface{}{
	"custom1": "#Custom1",
	"custom2": 4711,
	"custom3": []string{"Hallo", "Welt"},
}

var testObject = &A{C: 4711, Bs: []B{
	{Foo: "Hallo", Bar: 1, Baz: additionalData},
	{Foo: "Hallo2", Bar: 2, Baz: additionalData},
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

	if !reflect.DeepEqual(m, n) {
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
