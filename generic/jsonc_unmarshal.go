package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

func ObjectFromJson(j []byte, targetStruct interface{}) error {
	// Unmarshal json to the target struct
	err := json.Unmarshal(j, &targetStruct)
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

	return objectFromJson(targetStruct, additionalFieldsMap)
}

func objectFromJson(targetStruct interface{}, additionalFieldsMap map[string]interface{}) error {
	// is it a pointer of struct?
	structValue, ok := pointerOfStruct(&targetStruct)
	if ok == false {
		return errors.New("input is not a pointer of struct")
	}

	structType := structValue.Type()

	// Found field for "additional value" inside the struct
	var additionalFieldsName string

	// Iterate over all fields of the struct
	for i := 0; i < structType.NumField(); i++ {
		// Delete the struct fields from the additional fields
		delete(additionalFieldsMap, structType.Field(i).Name)

		// Get json fieldTag from the struct field and delete it from
		// the additionalFieldsMap
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)
		fieldTag := getJsonTag(&fieldType, "json")
		if fieldTag != nil {
			delete(additionalFieldsMap, fieldTag.Name)
		}

		// Find the jsonc field as additionalFieldsName
		fieldTag = getJsonTag(&fieldType, "jsonc")
		if fieldTag != nil {
			switch fieldTag.Name {
			case "flat":
				additionalFieldsName = fieldType.Name
				break
			case "collection":
				if fieldValue.Kind() != reflect.Slice {
					log.Printf("warn: Field %s ist not a slice!", fieldType.Name)
					break
				}

				println(fieldValue.Kind().String())
				break
			}
		}
	}

	// Add the additionalFieldsMap as attribute of the struct
	if additionalFieldsName != "" {
		field := structValue.FieldByName(additionalFieldsName)
		if field.Kind() == reflect.Map {
			field.Set(reflect.ValueOf(additionalFieldsMap))
		} else {
			log.Printf("Error: Field %s is not a map! Cannot deflat it.", additionalFieldsName)
		}
	}

	return nil
}
