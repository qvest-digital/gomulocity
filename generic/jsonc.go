package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
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
		field := valueType.Field(i)
		fieldName := field.Name

		// Ignore myself
		if fieldName == "JsonObject" {
			continue
		}

		// Handle special name "AdditionalFields" as map and flat it.
		if field.Tag == `jsonc:"flat"` {
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
			v, ok := field.Tag.Lookup("json")

			if ok {
				m[v] = value.Field(i).Interface()
			} else {
				m[fieldName] = value.Field(i).Interface()
			}
		}
	}

	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

func ObjectFromJson() {

}
