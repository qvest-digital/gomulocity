package generic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
)

func JsonFromObject(a interface{}) (string, error) {
	m := make(map[string]interface{})
	ptr := reflect.ValueOf(a)
	if ptr.Kind() != reflect.Ptr {
		return "", errors.New("No pointer of struct!")
	}

	value := ptr.Elem()
	if value.Kind() != reflect.Struct {
		return "", errors.New("No struct!")
	}

	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fmt.Println(t.Field(i).Tag)
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
			v := value.Field(i).Interface()
			m[fieldName] = v
		}
	}

	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}
