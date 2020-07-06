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
	tag := getJsonTag(fieldType)
	objectMap := *objectMapPtr

	if tag == nil {
		objectMap[fieldType.Name] = fieldValue.Interface()
		return
	}

	if tag.Name == "-" {
		return
	}

	if tag.OmitEmpty {
		if !isEmptyValue(fieldValue) {
			objectMap[tag.Name] = fieldValue.Interface()
		}
	} else {
		objectMap[tag.Name] = fieldValue.Interface()
	}
}

type Tag struct {
	Name      string
	OmitEmpty bool
}

func getJsonTag(fieldType *reflect.StructField) *Tag {
	tag, ok := fieldType.Tag.Lookup("json")
	if ok {
		tagValues := strings.Split(tag, ",")
		if len(tagValues) == 1 {
			return &Tag{tag, false}
		}

		if tagValues[1] == "omitempty" {
			return &Tag{tagValues[0], true}
		} else {
			return &Tag{tagValues[0], false}
		}
	}
	return nil
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

func ObjectFromJson(j []byte, target interface{}) error {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		if value.Kind() != reflect.Struct {
			return errors.New("input is not a pointer of struct")
		}
	} else {
		return errors.New("input is not a pointer of struct")
	}

	var tmpMap map[string]interface{}
	err := json.Unmarshal(j, &tmpMap)
	if err != nil {
		log.Printf("Error while unmarshaling json: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(j), &target)
	if err != nil {
		log.Printf("Error while unmarshaling json: %v", err)
		return err
	}

	// 1st: origin fields of the type
	valueType := value.Type()
	var typeFields = make([]string, valueType.NumField())
	for i := 0; i < valueType.NumField(); i++ {
		typeFields[i] = valueType.Field(i).Name
	}

	// 2snd: All tag renaming
	var typeTags = make(map[string]*Tag)
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		tag := getJsonTag(&field)
		if tag != nil {
			typeTags[field.Name] = tag
		}
	}

	// 3rd find flat field:
	var collectFieldName string
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		_, ok := field.Tag.Lookup("jsonc")
		if ok {
			collectFieldName = field.Name
		}
	}

	// Find all fields, that are not part of the original type
	for _, foo := range typeFields {
		delete(tmpMap, foo)
	}
	for _, foo := range typeTags {
		delete(tmpMap, foo.Name)
	}

	if collectFieldName != "" {
		value.FieldByName(collectFieldName).Set(reflect.ValueOf(tmpMap))
	}

	return err
}
