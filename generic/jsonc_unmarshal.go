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

	structType := structValue.Type()

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

	// Found field for "additional value" inside the struct
	var additionalFieldsName string

	// Iterate over all fields of the struct
	for i := 0; i < structType.NumField(); i++ {
		// Delete the struct fields from the additional fields
		delete(additionalFieldsMap, structType.Field(i).Name)

		// Get json fieldTag from the struct field and delete it from
		// the additionalFieldsMap
		field := structType.Field(i)
		fieldTag := getJsonTag(&field, "json")
		if fieldTag != nil {
			delete(additionalFieldsMap, fieldTag.Name)
		}

		// Find the jsonc field as additionalFieldsName
		fieldTag = getJsonTag(&field, "jsonc")
		if fieldTag != nil {
			switch fieldTag.Name {
			case "flat":
				additionalFieldsName = field.Name
				break
			case "collection":
				println(field.Name)
				break
			}
		}
	}

	for i := 0; i < structType.NumField(); i++ {

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

	return err
}
