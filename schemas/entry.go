package schemas

import "reflect"

type IEntry interface {
	GetTableName() string
}

type IEntries []IEntry

func GetPKsByEntries(entries interface{}) (PK, error) {
	pks := make(PK, 0)
	entriesElement := reflect.Indirect(reflect.ValueOf(entries))
	for i := 0; i < entriesElement.Len(); i++ {
		_, pk, err := GetEntryKey(reflect.Indirect(entriesElement.Index(i)).Interface().(IEntry))
		if err != nil {
			return pks, err
		}
		pks = append(pks, pk)
	}
	return pks, nil
}
