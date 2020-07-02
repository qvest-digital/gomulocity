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

func TestJsonc_ObjectFromJson_SuccessOnPlainStruct(t *testing.T) {
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
