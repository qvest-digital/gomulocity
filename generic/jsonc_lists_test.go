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
				"foo":"Hello",
				"bar":1,
				"custom1":"#Custom1",
				"custom2":4711,
				"custom3": [
					"Hallo",
					"Welt"
				]
			},
			{
				"foo":"Hello2",
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
	{Foo: "Hello", Bar: 1, Baz: additionalData},
	{Foo: "Hello2", Bar: 2, Baz: additionalData},
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

func TestJsonc_Marshal_Lists_WrongTags(t *testing.T) {
	type WrongFlat struct {
		A string `jsonc:"flat"`
	}
	type WrongCollection struct {
		A string `jsonc:"collection"`
	}

	_, err := JsonFromObject(&WrongFlat{A: "Hello"})
	if err == nil {
		t.Errorf("JsonFromObject - no error, want error for wrong use of jsonc:flat.")
	}

	_, err = JsonFromObject(&WrongCollection{A: "Hello"})
	if err == nil {
		t.Errorf("JsonFromObject - no error, want error for wrong use of jsonc:collection.")
	}
}

func TestJsonc_Unmarshal_Lists(t *testing.T) {
	a := &A{}
	err := ObjectFromJson([]byte(testJson), a)

	if err != nil {
		t.Errorf("ObjectFromJson - unexpected error %v", err)
	}

	if a.C != 4711 {
		t.Errorf("ObjectFromJson - basic elements = {C: %d}, want = {C: 4711}", a.C)
	}

	if len(a.Bs) != 2 {
		t.Errorf("ObjectFromJson - collection size = %d, want = 2", len(a.Bs))
	}

	assertB(a.Bs[0], "Hello", 1, t)
	assertB(a.Bs[1], "Hello2", 2, t)
}

func assertB(b B, foo string, bar int, t *testing.T) {
	if b.Bar != bar || b.Foo != foo {
		t.Errorf(
			"ObjectFromJson - basic elements = {Bar: %d, Foo: %s}, want = {Bar: %d, Foo: %s, FooBar: myFooBar}",
			b.Bar, b.Foo, bar, foo,
		)
	}

	custom1, _ := b.Baz["custom1"].(string)
	custom2, _ := b.Baz["custom2"].(float64)

	if custom1 != "#Custom1" || custom2 != 4711 {
		t.Errorf(
			"ObjectFromJson - Sub B = {custom1: %s, custom2: %.0f}, want = {custom1: %s, custom2: %f}",
			custom1, custom2, additionalData["custom1"], additionalData["custom2"],
		)
	}

	custom3, _ := b.Baz["custom3"].([]interface{})
	if len(custom3) != 2 {
		t.Errorf("ObjectFromJson - Sub B -> custom3 size = %d, want = 2", len(custom3))
	}
	if custom3[0] != "Hallo" || custom3[1] != "Welt" {
		t.Errorf("ObjectFromJson - Sub B -> custom3 = [%v, %v], want = [Hallo Welt]", custom3[0], custom3[1])
	}
}
