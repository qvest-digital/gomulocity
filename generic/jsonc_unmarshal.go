package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

/*
	Takes a json as []byte and a pointer of the target struct
	Returns an error, otherwise fills the `targetStruct` reference
	with values.
*/
func ObjectFromJson(j []byte, targetStruct interface{}) error {
	// is it a pointer of struct?
	structValue, ok := pointerOfStruct(&targetStruct)
	if ok == false {
		return errors.New("input is not a pointer of struct")
	}

	// First - let json unmarshal to the target struct as far as it gets
	err := json.Unmarshal(j, targetStruct)
	if err != nil {
		log.Printf("Error while unmarshaling json: %v", err)
		return err
	}

	// Second - Unmarshal json to a generic map: string -> interface
	// to have all data as raw fields as working structure.
	var fieldsMap map[string]interface{}
	err = json.Unmarshal(j, &fieldsMap)
	if err != nil {
		log.Printf("Error while unmarshaling json: %v", err)
		return err
	}

	return mergeMapWithStruct(&fieldsMap, structValue)
}

/*
	Merges the given object map into the struct.
	`structMapPtr` is a pointer of the object map
	`structValue` is the reflection Value of the struct object

	Returns an error, otherwise fills the `structValue` reference
	with values
*/
func mergeMapWithStruct(structMapPtr *map[string]interface{}, structValue *reflect.Value) error {
	structMap := *structMapPtr
	structType := structValue.Type()

	// Found field for "jsonc:"flat"" inside the struct
	var flatFieldName string

	// Iterate over all fields of the struct
	for i := 0; i < structType.NumField(); i++ {
		// Represents a single fields type and value
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Get tag values of the field
		jsonCTag := getJsonTag(&fieldType, "jsonc")
		jsonTag := getJsonTag(&fieldType, "json")

		if jsonCTag != nil {
			switch jsonCTag.Name {
			// The field is tagged with `jsonc:"flat"` -> set `flatFieldName`
			case "flat":
				if fieldValue.Kind() != reflect.Map {
					log.Printf("warn: Field %s is not a map!", fieldType.Name)
					break
				}

				flatFieldName = fieldType.Name
				break
			case "collection":
				// The field is tagged with `jsonc:"collection"`. Handle all elements as an flatted struct
				if fieldValue.Kind() != reflect.Slice {
					log.Printf("warn: Field %s ist not a slice!", fieldType.Name)
					break
				}

				// What is the json name of the collection in the `structMap`?
				var jsonFieldName string
				if jsonTag != nil {
					jsonFieldName = jsonTag.Name
				} else {
					jsonFieldName = fieldType.Name
				}

				// Get the collection from the `structMap`
				collection := structMap[jsonFieldName].([]interface{})

				// Iterate over each collection element
				for i := 0; i < fieldValue.Len(); i++ {
					// The collection element from the struct value and the struct map collection
					structElement := fieldValue.Index(i)
					mapElement, ok := collection[i].(map[string]interface{})
					if !ok {
						log.Printf("Element of collection field %s is not a map. Ignoring it!", fieldType.Name)
						break
					}

					// Call this function recursively with the collection element.
					err := mergeMapWithStruct(&mapElement, &structElement)
					if err != nil {
						log.Printf("error while unmarshaling colletion field %s", fieldType.Name)
					}
				}
				break
			}
		}

		// At the end, `structMap` must contain only the "non struct" fields.
		// Therefore, delete all "known" struct fields from the `structMap`
		if jsonTag != nil {
			delete(structMap, jsonTag.Name)
		}
		delete(structMap, structType.Field(i).Name)
	}

	// Add the structMap as value of the struct field `flatFieldName` or `jsonc:"flat"`
	if flatFieldName != "" {
		field := structValue.FieldByName(flatFieldName)
		field.Set(reflect.ValueOf(structMap))
	}

	return nil
}
