package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

func TrimSpacesFromStruct(data interface{}) {
	val := reflect.ValueOf(data).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldType := field.Type
		fieldName := field.Tag.Get("json")
		if fieldType.Kind() == reflect.String && fieldName != "password" {
			fieldValue := val.Field(i).Interface().(string)
			val.Field(i).SetString(strings.TrimSpace(fieldValue))
		}
	}
}

func PrintStructFields(data interface{}) {
	val := reflect.ValueOf(data)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Name
		fieldValue := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", fieldName, fieldValue)
	}
}
