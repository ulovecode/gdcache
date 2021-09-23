package reflect

import (
	"reflect"
)

// IsPointerElementSlice Is a pointer type item
func IsPointerElementSlice(value interface{}) bool {
	typeValue := reflect.Indirect(reflect.ValueOf(value)).Type().Elem()
	return typeValue.Kind() == reflect.Ptr
}

// GetSliceValue Get the slice type value of the interface
func GetSliceValue(value interface{}) interface{} {
	var v reflect.Value
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		v = reflect.ValueOf(value).Elem()
	} else {
		v = reflect.ValueOf(value)
	}
	typeValue := v.Type().Elem()
	return reflect.New(typeValue).Interface()
}

// CovertSlicePointerValue2StructValue Convert slice type pointer array to slice type entity array
func CovertSlicePointerValue2StructValue(sliceInterface interface{}) interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(sliceInterface))
	newSlice := reflect.MakeSlice(reflect.SliceOf(sliceValue.Type().Elem().Elem()), sliceValue.Len(), sliceValue.Len())
	for i := 0; i < newSlice.Len(); i++ {
		newSlice.Index(i).Set(reflect.Indirect(sliceValue.Index(i)))
	}
	return newSlice.Interface()
}

// CovertSliceStructValue2PointerValue Convert a slice type entity array to a slice type pointer array
func CovertSliceStructValue2PointerValue(sliceInterface interface{}) interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(sliceInterface))

	newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(sliceValue.Type().Elem())), sliceValue.Len(), sliceValue.Len())
	for i := 0; i < newSlice.Len(); i++ {
		newSlice.Index(i).Set(reflect.Indirect(sliceValue.Index(i)).Addr())
	}
	return newSlice.Interface()
}

// MakePointerSliceValue Create a pointer type slice
func MakePointerSliceValue(entriesValue reflect.Value) reflect.Value {
	slice := reflect.MakeSlice(reflect.Indirect(entriesValue).Type(), 0, 0)
	value := reflect.New(slice.Type())
	value.Elem().Set(slice)
	return value
}
