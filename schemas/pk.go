package schemas

import "strings"

// PK represents primary key values
type PK []string

// ToEntryKeys Convert the cache key to the corresponding entity key
func (pk PK) ToEntryKeys() []EntryKeys {
	entryKeys := make([]EntryKeys, 0)
	for _, p := range pk {
		entryKey := make([]EntryKey, 0)
		entryKeyString := strings.SplitN(p, "#", 2)
		keyValues := strings.Split(entryKeyString[1], "-")
		for _, keyValue := range keyValues {
			keyValuesString := keyValue[1 : len(keyValue)-1]
			keyValueArray := strings.Split(keyValuesString, ":")
			key := keyValueArray[0]
			value := keyValueArray[1]
			entryKey = append(entryKey, EntryKey{
				Name:  key,
				Param: value,
			})
		}
		entryKeys = append(entryKeys, entryKey)
	}
	return entryKeys
}
