package generic

import (
	"testing"
)

func TestJsonc_Lists(t *testing.T) {
	type B struct {
		Foo string                 `json:"foo"`
		Bar int                    `json:"bar"`
		Baz map[string]interface{} `jsonc:"flat"`
	}

	type A struct {
		Bs []B `json:"bList"`
		C  int `json:"c"`
	}

	additional := map[string]interface{}{
		"custom1": "#Custom1",
		"custom2": 4711,
		"custom3": []string{"Hallo", "Welt"},
	}
	a := &A{C: 4711, Bs: []B{
		{Foo: "Hallo", Bar: 1, Baz: additional},
		{Foo: "Hallo2", Bar: 2, Baz: additional},
	}}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{
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
	if j != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}
