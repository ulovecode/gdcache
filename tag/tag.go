package tag

import (
	"reflect"
)

var defaultTag = &tag{tagName: "cache"}

type tag struct {
	tagName string
}

// ConfigTag Configure cache tag name
func ConfigTag(tagName string) {
	defaultTag = &tag{tagName: tagName}
}

// GetName Get tag name
func GetName() string {
	return defaultTag.tagName
}

// GetCacheTagFields Get the cached tag field on the entity
func GetCacheTagFields(value interface{}) []reflect.StructField {
	reflectValue := reflect.TypeOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	structFields := make([]reflect.StructField, 0)
	for i := 0; i < reflectValue.NumField(); i++ {
		structTag := reflectValue.Field(i).Tag
		if _, ok := structTag.Lookup(defaultTag.tagName); ok {
			structFields = append(structFields, reflectValue.Field(i))
		}
	}
	return structFields
}
