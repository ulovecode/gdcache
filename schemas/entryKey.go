package schemas

import (
	"errors"
	"fmt"
	"github.com/ulovecode/gdcache/tag"
	"reflect"
	"strings"
)

var ServiceName string

type EntryKey struct {
	Name  string
	Param string
}

type EntryKeys []EntryKey

func (es EntryKeys) GetEntryKey(entryName string) string {
	var (
		keyTemplate   = make([]string, 0)
		entryKeyNames = make([]interface{}, 0)
	)
	for _, e := range es {
		keyTemplate = append(keyTemplate, fmt.Sprintf("[%s", e.Name)+":%s]")
		entryKeyNames = append(entryKeyNames, e.Param)
	}
	return fmt.Sprintf(ServiceName+"_"+entryName+"#"+strings.Join(keyTemplate, "-"), entryKeyNames...)
}

// GetEntryKey get the cache primary Name, if not, find the default value field as id
func GetEntryKey(entry IEntry) ([]EntryKey, string, error) {

	var (
		entryKeys  EntryKeys = make([]EntryKey, 0)
		entryValue           = reflect.ValueOf(entry)
	)

	switch entryValue.Type().Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		entryValue = entryValue.Elem()
	}

	tagSortFields := tag.GetCacheTagFields(entry)

	if len(tagSortFields) == 0 {
		fieldValue := entryValue.FieldByNameFunc(func(fileName string) bool {
			if strings.ToLower(fileName) == "id" {
				return true
			}
			return false
		})
		if fieldValue == reflect.ValueOf(nil) {
			return entryKeys, "", errors.New("The field with the default value of Id was not found, and the setting cache tag was not found either ")
		}
		param := fieldValue.Interface()
		entryKeys = append(entryKeys, EntryKey{
			Name:  "id",
			Param: fmt.Sprint(param),
		})
		return entryKeys, entryKeys.GetEntryKey(entryValue.Type().String()), nil
	}

	for _, tagField := range tagSortFields {
		fieldValue := entryValue.FieldByIndex(tagField.Index)
		entryKeys = append(entryKeys, EntryKey{
			Name:  tagField.Tag.Get(tag.GetName()),
			Param: fmt.Sprint(fieldValue.Interface()),
		})
	}
	return entryKeys, entryKeys.GetEntryKey(entryValue.Type().String()), nil
}
