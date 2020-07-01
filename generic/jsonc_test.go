package generic

import (
	"testing"
)

func TestJsonc_ErrorOnNoStruct(t *testing.T) {
	j, err := JsonFromObject("Hallo Welt")

	if err == nil {
		t.Errorf("JsonFromObject - error expected. Instead: %v", j)
	}
}

func TestJsonc_SuccessOnPlainStruct(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := A{B: "Foo", C: 4711}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"B":"Foo","C":4711}`
	if j != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_SuccessOnPointerStruct(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := &A{B: "Foo", C: 4711}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"B":"Foo","C":4711}`
	if j != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_MarshalStandardFields(t *testing.T) {
	type A struct {
		B string
		C int
	}

	a := &A{B: "Foo", C: 4711}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"B":"Foo","C":4711}`
	if j != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_SupportsStandardJsonTags(t *testing.T) {
	type A struct {
		B string `json:"myB"`
		C int    `json:"myC"`
	}

	a := &A{B: "Foo", C: 4711}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"myB":"Foo","myC":4711}`
	if j != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}

func TestJsonc_DoesNotFlatUntaggedMaps(t *testing.T) {
	type A struct {
		B string
		C map[string]string
	}

	a := &A{B: "Foo", C: map[string]string{"foo1": "bar", "foo2": "baz"}}

	j, err := JsonFromObject(a)

	if err != nil {
		t.Errorf("JsonFromObject - unexpected error %v", err)
	}

	want := `{"B":"Foo","C":{"foo1":"bar","foo2":"baz"}}`
	if j != want {
		t.Errorf("JsonFromObject - json = %v, want %v", j, want)
	}
}
