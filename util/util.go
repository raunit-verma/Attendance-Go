package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const TimeZone = "Asia/Kolkata"

// checks if string is valid email or not
func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

// removes spaces from struct expect from password field
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

// Prints any struct with field name and value
func PrintStructFields(data interface{}) {
	val := reflect.ValueOf(data)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Name
		fieldValue := val.Field(i).Interface()
		fmt.Printf("%s: %v\n", fieldName, fieldValue)
	}
}

// Converts and return current time in India timezone
func GetCurrentIndianTime() time.Time {
	currentTimeUTC := time.Now().UTC()
	indiaTimeZone, err := time.LoadLocation(TimeZone)
	if err != nil {
		zap.L().Error("Error loading indian timezone")
		return time.Now()
	}
	return currentTimeUTC.In(indiaTimeZone)
}

// Formate Date Time to Required Format
func FormateDateTime(year int, month time.Month, date int, hour int, min int, sec int) (string, time.Time) {
	punchInDate := time.Date(year, month, date, hour, min, sec, 0, time.UTC)
	indiaTimeZone, err := time.LoadLocation(TimeZone)
	if err != nil {
		zap.L().Error("Error loading time zone")
		return "", time.Now()
	}
	punchInDateIST := punchInDate.In(indiaTimeZone)
	layout := "2006-01-02 15:04:05-07:00"
	timeString := punchInDateIST.Format(layout)
	parsedTime, err := time.Parse(layout, timeString)
	return timeString, parsedTime
}

// Match hashed password with bcrypt
func MatchPassword(hashedPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return false
	}
	return true
}

// Hash Password
func GenerateHashFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Error in creating hash from password.", zap.Error(err))
		return "", err
	}
	return string(hashedPassword), nil
}

// strong password checker
func IsStrongPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, " Password must be of 8 characters."
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpperCase && hasLowerCase && hasDigit, " Password must have one uppercase, one lowercase & one digit."
}
