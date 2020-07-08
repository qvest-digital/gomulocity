package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

func ObjectFromJson(j []byte, targetStruct interface{}) error {
	// is it a pointer of struct?
	structValue, ok := pointerOfStruct(&targetStruct)
	if ok == false {
		return errors.New("input is not a pointer of struct")
	}

	// Unmarshal json to the target struct
	err := json.Unmarshal(j, targetStruct)
	if err != nil {
		log.Printf("Error while unmarshaling json: %v", err)
		return err
	}

	// Unmarshal json to a generic map: string -> interface
	var additionalFieldsMap map[string]interface{}
	err = json.Unmarshal(j, &additionalFieldsMap)
	if err != nil {
		log.Printf("Error while unmarshaling json: %v", err)
		return err
	}

	return objectFromJson(structValue, &additionalFieldsMap)
}

func objectFromJson(structValue *reflect.Value, additionalFieldsMapPtr *map[string]interface{}) error {
	additionalFieldsMap := *additionalFieldsMapPtr
	structType := structValue.Type()

	// Found field for "additional value" inside the struct
	var additionalFieldsName string

	// Iterate over all fields of the struct
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Find the jsonc field as additionalFieldsName
		jsonCTag := getJsonTag(&fieldType, "jsonc")
		jsonTag := getJsonTag(&fieldType, "json")

		if jsonCTag != nil {
			switch jsonCTag.Name {
			case "flat":
				if fieldValue.Kind() != reflect.Map {
					log.Printf("warn: Field %s is not a map!.", fieldType.Name)
					break
				}

				additionalFieldsName = fieldType.Name
				break
			case "collection":
				if fieldValue.Kind() != reflect.Slice {
					log.Printf("warn: Field %s ist not a slice!", fieldType.Name)
					break
				}

				var sliceFieldName string
				if jsonTag != nil {
					sliceFieldName = jsonTag.Name
				} else {
					sliceFieldName = fieldType.Name
				}
				bar := additionalFieldsMap[sliceFieldName].([]interface{})

				for i := 0; i < fieldValue.Len(); i++ {
					o := fieldValue.Index(i)
					w := bar[i].(map[string]interface{})
					println(w)
					err := objectFromJson(&o, &w)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
				break
			}
		}

		if jsonTag != nil {
			delete(additionalFieldsMap, jsonTag.Name)
		}

		// Delete the struct fields from the additional fields
		delete(additionalFieldsMap, structType.Field(i).Name)

		// Get json fieldTag from the struct field and delete it from
		// the additionalFieldsMap
	}

	// Add the additionalFieldsMap as attribute of the struct
	if additionalFieldsName != "" {
		field := structValue.FieldByName(additionalFieldsName)
		field.Set(reflect.ValueOf(additionalFieldsMap))
	}

	return nil
}
