package generic

import (
	"reflect"
	"testing"
)

func TestJsonc_ObjectFromJson_ErrorOnInvalidJson(t *testing.T) {
	// given: A simple struct
	type A struct{}
	a := &A{}

	// when: We unmarshal with an invalid string
	err := ObjectFromJson([]byte("Hallo Welt"), a)

	// then:
	if err == nil {
		t.Errorf("ObjectFromJson - error expected. Instead: %v", a)
	}
}

func TestJsonc_ObjectFromJson_ErrorOnPlainStruct(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := A{}

	err := ObjectFromJson([]byte(`{"B": "Hello", "C": 4711 }`), a)

	if err == nil {
		t.Errorf("ObjectFromJson - error expected. Instead: %v", a)
	}
}

func TestJsonc_ObjectFromJson_SuccessOnPointerStruct(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := &A{}

	err := ObjectFromJson([]byte(`{"B": "Hello", "C": 4711 }`), a)

	if err != nil {
		t.Errorf("ObjectFromJson - unexpected error %v", err)
	}

	want := &A{B: "Hello", C: 4711}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("ObjectFromJson - object = %v, want %v", a, want)
	}
}

func TestJsonc_ObjectFromJson_SupportsStandardJsonTags(t *testing.T) {
	type A struct {
		B string `json:"myB"`
		C int    `json:"myC"`
		D string `json:"-"`
		E string `json:"myE,omitempty"`
		F bool   `json:"myF,omitempty"`
		G int    `json:",omitempty"`
		H *A     `json:"myH,omitempty"`
		I string `json:"myI,otherstuff"`
	}

	a := &A{}
	j := `{"myB":"Foo", "myC":4711, "D":"Bar", "myE": null, "myF": false, "G": 0, "myH": null, "myI": "Hello"}`
	err := ObjectFromJson([]byte(j), a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := &A{B: "Foo", C: 4711, E: "", F: false, G: 0, H: nil, I: "Hello"}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("ObjectFromJson - object = %v, want %v", a, want)
	}
}

func TestJsonc_ObjectFromJson_CollectOtherFieldsInD(t *testing.T) {
	// given: A struct with field B and C as well defined types and a field D as generic bucket.
	type A struct {
		B string                 `json:"myB"`
		C int                    `json:"myC"`
		D map[string]interface{} `jsonc:"flat"`
	}

	// and: A test json, with the fields B, C and other fields vom E to I.
	a := &A{}
	j := `{"myB":"Foo", "myC":4711, "myE": null, "myF": false, "G": 0, "myH": 0.567, "myI": [ "Hello", "Welt" ]}`

	// when: We unmarshal the json
	err := ObjectFromJson([]byte(j), a)

	// then: We do no expect an error
	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	// and: B and C hat correct data
	if a.B != "Foo" || a.C != 4711 {
		t.Errorf("ObjectFromJson - basic elements = {B: %s, C: %d}, want = {B: %s, C: %d}", a.B, a.C, "Foo", 4711)
	}

	// and: E to H has correct types and data
	myE, _ := a.D["myE"].(interface{})
	myF, _ := a.D["myF"].(bool)
	myG, _ := a.D["myG"].(int)
	myH, _ := a.D["myH"].(float64)

	if myE != nil ||
		myF != false ||
		myG != 0 ||
		myH != 0.567 {
		t.Errorf("ObjectFromJson - D = {%v, %v, %v, %v}, want {%v, %v, %v, %v}", myE, myF, myG, myH, nil, false, 0, 0.567)
	}

	// and: I is a slice with correct data
	myI, _ := a.D["myI"].([]interface{})
	hello, _ := myI[0].(string)
	world, _ := myI[1].(string)

	if hello != "Hello" && world != "World" {
		t.Errorf("ObjectFromJson - D.myI = [ %v, %v ], want [ %v, %v ]", hello, world, "Hello", "World")
	}
}
