package reflect

import (
	"github.com/ulovecode/gdcache/schemas"
	"reflect"
)

func GetSliceValue(value interface{}) schemas.IEntry {
	typeValue := reflect.TypeOf(value).Elem().Elem()
	return reflect.New(typeValue).Interface().(schemas.IEntry)
}

func MakeEmptySlice(value interface{}) []schemas.IEntry {
	return reflect.New(reflect.TypeOf(value)).Interface().([]schemas.IEntry)
}
