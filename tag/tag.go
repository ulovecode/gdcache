package tag

import (
	"reflect"
	"sort"
	"strings"
)

var defaultTag *tag

type tag struct {
	tagName string
}

// ConfigTag Configure cache tag name
func ConfigTag(tagName string) {
	defaultTag = &tag{tagName: tagName}
}

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
		if _, ok := reflectValue.Field(i).Tag.Lookup(defaultTag.tagName); ok {
			structFields = append(structFields, reflectValue.Field(i))
		}
	}
	sort.Slice(structFields, func(i, j int) bool {
		iTag := strings.TrimSpace(structFields[i].Tag.Get(defaultTag.tagName))
		jTag := strings.TrimSpace(structFields[j].Tag.Get(defaultTag.tagName))
		if len(iTag) < len(jTag) {
			return true
		}
		if jTag < jTag {
			return true
		}
		return false
	})
	return structFields
}
