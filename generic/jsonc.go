package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
)

func JsonFromObject(a interface{}) (string, error) {
	value := reflect.ValueOf(a)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return "", errors.New("input is not a struct or pointer of struct")
	}

	m := make(map[string]interface{})
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		fieldType := valueType.Field(i)
		fieldName := fieldType.Name

		// Ignore myself
		if fieldName == "JsonObject" {
			continue
		}

		// Handle special name "jsonc:"flat"" as map and flat it.
		if fieldType.Tag == `jsonc:"flat"` {
			fieldValue := value.Field(i)
			if fieldValue.Kind() != reflect.Map {
				log.Printf("error: %s is not a map. Can not flat it into the json and therefore ignore it.", fieldName)
				continue
			}

			iter := fieldValue.MapRange()
			for iter.Next() {
				m[iter.Key().String()] = iter.Value().Interface()
			}
		} else {
			field := value.Field(i)
			insertIntoMap(&m, &fieldType, &field)
		}
	}

	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

func insertIntoMap(objectMapPtr *map[string]interface{}, fieldType *reflect.StructField, fieldValue *reflect.Value) {
	tag, ok := fieldType.Tag.Lookup("json")
	objectMap := *objectMapPtr

	if !ok {
		objectMap[fieldType.Name] = fieldValue.Interface()
		return
	}

	if tag == "-" {
		return
	}

	tagValues := strings.Split(tag, ",")
	if len(tagValues) == 1 {
		objectMap[tag] = fieldValue.Interface()
		return
	}

	if tagValues[1] == "omitempty" {
		if !isEmptyValue(fieldValue) {
			objectMap[tag] = fieldValue.Interface()
		}
	} else {
		objectMap[tagValues[0]] = fieldValue.Interface()
	}

	return
}

func isEmptyValue(v *reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func ObjectFromJson() {

}
