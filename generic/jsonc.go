package generic

import (
	"reflect"
	"strings"
)

type Tag struct {
	Type      string
	Field     string
	Name      string
	OmitEmpty bool
}

/**
Takes a possible pointer on struct `value *reflect.Value`
Returns true/false, whether it is a pointer on struct or not.
Returns the *reflect.Value representing the struct.
*/
func pointerOfStruct(o *interface{}) (*reflect.Value, bool) {
	value := reflect.ValueOf(*o)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		if value.Kind() == reflect.Struct {
			return &value, true
		} else {
			return nil, false
		}
	} else {
		return nil, false
	}
}

func getJsonTag(fieldType *reflect.StructField, tagName string) *Tag {
	tag, ok := fieldType.Tag.Lookup(tagName)
	if !ok {
		return nil
	}

	tagValues := strings.Split(tag, ",")
	if len(tagValues) == 2 && tagValues[1] == "omitempty" {
		return &Tag{tagName, fieldType.Name, tagValues[0], true}
	} else {
		return &Tag{tagName, fieldType.Name, tagValues[0], false}
	}
}
