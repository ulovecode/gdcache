package reflect

import (
	"gdcache/schemas"
	"reflect"
)

func GetSliceValue(value interface{}) schemas.IEntry {
	typeValue := reflect.TypeOf(value).Elem().Elem()
	//fmt.Printf("Test: %v |||",typeValue.String() )
	return reflect.New(typeValue).Interface().(schemas.IEntry)
}

func MakeEmptySlice(value interface{}) []schemas.IEntry {
	return reflect.New(reflect.TypeOf(value)).Interface().([]schemas.IEntry)
}
