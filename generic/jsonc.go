package generic

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
)

const flatTag = `jsonc:"flat"`

type Tag struct {
	Name      string
	OmitEmpty bool
}

func JsonFromObject(o interface{}) (string, error) {
	value, ok := pointerOfStruct(&o)
	if ok == false {
		return "", errors.New("input is not a pointer of struct")
	}

	m := make(map[string]interface{})
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		fieldType := valueType.Field(i)
		fieldName := fieldType.Name

		// Handle special name "jsonc:"flat"" as map and flat it.
		if fieldType.Tag == flatTag {
			fieldValue := value.Field(i)
			if fieldValue.Kind() != reflect.Map {
				log.Printf("error: %s is not a map. Can not flat it into the json and therefore ignore it.", fieldName)
				continue
			}

			// flat process
			iter := fieldValue.MapRange()
			for iter.Next() {
				m[iter.Key().String()] = iter.Value().Interface()
			}
		} else {
			field := value.Field(i)
			insertTaggedFieldIntoMap(&m, &fieldType, &field)
		}
	}

	j, err := json.Marshal(m)
	if err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

func ObjectFromJson(j []byte, targetStruct interface{}) error {
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
		tag := getJsonTag(&field)
		if tag != nil {
			delete(additionalFieldsMap, tag.Name)
		}

		// Find the jsonc field as collectFieldName
		_, ok := field.Tag.Lookup("jsonc")
		if ok {
			collectFieldName = field.Name
		}
	}

	// The the additionalFieldsMap as attribute of the struct
	if collectFieldName != "" {
		value.FieldByName(collectFieldName).Set(reflect.ValueOf(additionalFieldsMap))
	}

	return err
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

func insertTaggedFieldIntoMap(objectMapPtr *map[string]interface{}, fieldType *reflect.StructField, fieldValue *reflect.Value) {
	tag := getJsonTag(fieldType)
	objectMap := *objectMapPtr

	// no tag -> original name
	if tag == nil {
		objectMap[fieldType.Name] = fieldValue.Interface()
		return
	}

	// - -> omit value
	if tag.Name == "-" {
		return
	}

	// OmitEmpty and is empty -> omit value
	if tag.OmitEmpty && isEmptyValue(fieldValue) {
		return
	} else {
		objectMap[tag.Name] = fieldValue.Interface()
	}
}

func getJsonTag(fieldType *reflect.StructField) *Tag {
	tag, ok := fieldType.Tag.Lookup("json")
	if !ok {
		return nil
	}

	tagValues := strings.Split(tag, ",")
	if len(tagValues) == 2 && tagValues[1] == "omitempty" {
		return &Tag{tagValues[0], true}
	} else {
		return &Tag{tagValues[0], false}
	}
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
