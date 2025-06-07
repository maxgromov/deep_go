package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	v := reflect.ValueOf(person)
	vType := v.Type()
	res := make([]string, 0, vType.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := vType.Field(i)

		key, options := parseTag(field)
		if key == "" {
			continue
		}

		fieldValue := v.Field(i)
		if shouldSkipField(options, fieldValue) {
			continue
		}

		valueStr := formatValue(fieldValue)
		res = append(res, key+"="+valueStr)

	}

	return strings.Join(res, "\n")
}

func parseTag(field reflect.StructField) (key string, options []string) {
	tag := field.Tag.Get("properties")
	tagParts := strings.Split(tag, ",")

	return tagParts[0], tagParts[1:]
}

func formatValue(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Bool:
		return fmt.Sprintf("%t", val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}

func shouldSkipField(opts []string, val reflect.Value) bool {
	for _, opt := range opts {
		if opt == "omitempty" && val.IsZero() {
			return true
		}
	}
	return false
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
