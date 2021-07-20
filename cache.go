package gdcache

import (
	"errors"
	"fmt"
	"gdcache/schemas"
	"gdcache/tags"
	"reflect"
	"strings"
)

type ICache interface {
	Store(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type ICacheHandler interface {
	// GetEntries cache the entity content obtained through sql, and return the entity of the array pointer type
	GetEntries(entry interface{}, sql string, option ...Options) (interface{}, error)
	// GetEntry get a pointer to an entity type and return the entity
	GetEntry(entry interface{}, option ...Options) (interface{}, bool, error)
	// PutEntry trigger cache update through query results
	PutEntry(entry interface{}, sql string, option ...Options) error
	// DelEntries delete the corresponding execution entity through sql,
	// Before the update, you can query the id to be deleted first, and then delete these through defer
	DelEntries(entry interface{}, ids []interface{}, option ...Options) error
}

type CacheHandler struct {
	cacheHandler    ICache
	databaseHandler IDB
	encoder         IEncoder
	serializer      Serializer
	tags            tags.ITag
}

func (c CacheHandler) PutEntry(entry interface{}, sql string, option ...Options) error {
	panic("implement me")
}

func (c CacheHandler) GetEntries(entry interface{}, sql string, option ...Options) (interface{}, error) {
	panic("implement me")
}

func (c CacheHandler) GetEntry(entry interface{}, option ...Options) (interface{}, bool, error) {
	panic("implement me")
}

func (c CacheHandler) DelEntries(entry interface{}, ids []interface{}, option ...Options) error {
	panic("implement me")
}

func (c CacheHandler) getEntriesBySQL(entries interface{}, sql string) error {
	err := c.databaseHandler.Exec(sql, entries)
	return err
}

func (c CacheHandler) getIdsByEntries(entries interface{}) ([]schemas.PK, error) {
	entriesValue := reflect.ValueOf(entries)
	pks := make([]schemas.PK, 0)
	for i := 0; i < entriesValue.Len(); i++ {
		pk, err := c.getEntryKey(entriesValue.Index(i).Interface())
		if err != nil {
			return pks, err
		}
		pks = append(pks, pk)
	}
	return pks, nil
}

func (c CacheHandler) getIdsByCacheSQL(sql string) ([]schemas.PK, error) {
	pks := make([]schemas.PK, 0)
	ids, err := c.cacheHandler.Get(sql)
	if err != nil {
		return pks, err
	}
	return c.encoder.decodeIds(ids)
}

func (c CacheHandler) storeCacheSQLToIds(sql string) ([]schemas.PK, error) {
	pks := make([]schemas.PK, 0)
	pks, err := c.getIdsByCacheSQL(sql)
	if err != nil {
		return pks, err
	}
	bytes, err := c.serializer.serialize(pks)
	if err != nil {
		return pks, err
	}
	err = c.cacheHandler.Store(sql, bytes)
	return pks, err
}

// getEntryKey get the cache primary key, if not, find the default value field as id
func (c CacheHandler) getEntryKey(entry interface{}) (string, error) {
	tagSortFields := c.tags.GetCacheTagSortFields(entry)
	entryValue := reflect.ValueOf(entry)
	if len(tagSortFields) == 0 {
		fieldValue := entryValue.FieldByNameFunc(func(fileName string) bool {
			if fileName == "id" {
				return true
			}
			return false
		})
		if fieldValue == reflect.ValueOf(nil) {
			return "", errors.New("The field with the default value of Id was not found, and the setting cache tag was not found either ")
		}
		return fmt.Sprintf("[id:%s]", fieldValue.Interface()), nil
	}
	var keyTemplate = make([]string, 0)
	entryKeys := make([]interface{}, 0)
	for _, tagField := range tagSortFields {
		fieldValue := entryValue.FieldByIndex(tagField.Index)

		keyTemplate = append(keyTemplate, fmt.Sprintf("[%s", tagField.Name)+":%s]")
		entryKeys = append(entryKeys, fieldValue.Interface())
	}
	return fmt.Sprintf(strings.Join(keyTemplate, ":"), entryKeys...), nil
}
