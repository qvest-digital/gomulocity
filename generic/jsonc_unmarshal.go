package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

func ObjectFromJson(j []byte, targetStruct interface{}) error {
	// is it a pointer of struct?
	value, ok := pointerOfStruct(&targetStruct)
	if ok == false {
		return errors.New("input is not a pointer of struct")
	}

	valueType := value.Type()

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
	var collectFieldName string

	// Iterate over all fields of the struct
	for i := 0; i < valueType.NumField(); i++ {
		// Delete the struct fields from the additional fields
		delete(additionalFieldsMap, valueType.Field(i).Name)

		// Get json tag from the struct field and delete it from
		// the additionalFieldsMap
		field := valueType.Field(i)
		tag := getJsonTag(&field, "json")
		if tag != nil {
			delete(additionalFieldsMap, tag.Name)
		}

		// Find the jsonc field as collectFieldName
		tag = getJsonTag(&field, "jsonc")
		if tag != nil {
			if tag.Name == "flat" {
				collectFieldName = field.Name
			}
		}
	}

	// The the additionalFieldsMap as attribute of the struct
	if collectFieldName != "" {
		value.FieldByName(collectFieldName).Set(reflect.ValueOf(additionalFieldsMap))
	}

	return err
}
