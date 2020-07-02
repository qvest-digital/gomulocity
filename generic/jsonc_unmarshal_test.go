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
	err := ObjectFromJson("Hallo Welt", a)

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

	err := ObjectFromJson(`{"B": "Hello", "C": 4711 }`, a)

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

	err := ObjectFromJson(`{"B": "Hello", "C": 4711 }`, a)

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
	err := ObjectFromJson(j, a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := &A{B: "Foo", C: 4711, E: "", F: false, G: 0, H: nil, I: "Hello"}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("ObjectFromJson - object = %v, want %v", a, want)
	}
}
