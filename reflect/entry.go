package reflect

import (
	"reflect"
)

func IsPointerElement(value interface{}) bool {
	typeValue := reflect.Indirect(reflect.ValueOf(value)).Type().Elem()
	return typeValue.Kind() == reflect.Ptr
}

func GetSliceValue(value interface{}) interface{} {
	var v reflect.Value
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		v = reflect.ValueOf(value).Elem()
	} else {
		v = reflect.ValueOf(value)
	}
	typeValue := reflect.TypeOf(v)
	return reflect.New(typeValue).Interface()
}

func CovertSlicePointerValue2StructValue(sliceInterface interface{}) interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(sliceInterface))
	newSlice := reflect.MakeSlice(reflect.SliceOf(sliceValue.Type().Elem().Elem()), sliceValue.Len(), sliceValue.Len())
	for i := 0; i < newSlice.Len(); i++ {
		newSlice.Index(i).Set(reflect.Indirect(sliceValue.Index(i)))
	}
	return newSlice.Interface()
}

func CovertSliceStructValue2PointerValue(sliceInterface interface{}) interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(sliceInterface))

	newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(sliceValue.Type().Elem())), sliceValue.Len(), sliceValue.Len())
	for i := 0; i < newSlice.Len(); i++ {
		newSlice.Index(i).Set(reflect.Indirect(sliceValue.Index(i)).Addr())
	}
	return newSlice.Interface()
}
