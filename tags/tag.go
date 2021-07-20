package tags

import (
	"reflect"
	"sort"
	"strings"
)

type ITag interface {
	GetCacheTagSortFields(value interface{}) []reflect.StructField
}

type Tag struct {
	tagName string
}

var _ ITag = Tag{}

func (t Tag) GetCacheTagSortFields(value interface{}) []reflect.StructField {
	reflectValue := reflect.TypeOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	structFields := make([]reflect.StructField, 0)
	for i := 0; i < reflectValue.NumField(); i++ {
		if _, ok := reflectValue.Field(i).Tag.Lookup(t.tagName); ok {
			structFields = append(structFields, reflectValue.Field(i))
		}
	}
	sort.Slice(structFields, func(i, j int) bool {
		iTag := strings.TrimSpace(structFields[i].Tag.Get(t.tagName))
		jTag := strings.TrimSpace(structFields[j].Tag.Get(t.tagName))
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
