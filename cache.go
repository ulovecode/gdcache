package gdcache

import (
	"fmt"
	"github.com/ulovecode/gdcache/builder"
	"github.com/ulovecode/gdcache/log"
	gdreflect "github.com/ulovecode/gdcache/reflect"
	"github.com/ulovecode/gdcache/schemas"
	"github.com/ulovecode/gdcache/tag"
	"reflect"
)

type ReturnKeyValue struct {
	KeyValue
	Has bool
}

type KeyValue struct {
	Key   string
	Value []byte
}

type ICache interface {
	StoreAll(keyValues ...KeyValue) (err error)
	Get(key string) (data []byte, has bool, err error)
	GetAll(keys schemas.PK) (data []ReturnKeyValue, err error)
	DeleteAll(keys schemas.PK) error
}

type ICacheHandler interface {
	// GetEntries cache the entity content obtained through sql, and return the entity of the array pointer type
	GetEntries(entrySlice interface{}, sql string) error
	// GetEntry get a pointer to an entity type and return the entity
	GetEntry(entry interface{}) (bool, error)
	// DelEntries delete the corresponding execution entity through sql,
	// Before the update, you can query the id to be deleted first, and then delete these through defer
	DelEntries(entrySlice interface{}, sql string) error
}

var _ ICacheHandler = CacheHandler{}

type CacheHandler struct {
	cacheHandler    ICache
	databaseHandler IDB
	serializer      Serializer
	log             log.Logger
}

func NewCacheHandler(cacheHandler ICache, databaseHandler IDB, options ...OptionsFunc) *CacheHandler {
	o := Options{}
	for _, option := range options {
		option(&o)
	}

	tag.ConfigTag(o.cacheTagName)
	if tag.GetName() == "" {
		tag.ConfigTag("cache")
	}

	if o.serializer == nil {
		o.serializer = JsonSerializer{}
	}

	if o.log == nil {
		o.log = log.DefaultLogger{}
	}

	return &CacheHandler{cacheHandler: cacheHandler, databaseHandler: databaseHandler, serializer: o.serializer, log: o.log}
}

func (c CacheHandler) GetEntries(entrySlice interface{}, sql string) error {

	entriesValue := reflect.Indirect(reflect.ValueOf(entrySlice))
	entryElementType := entriesValue.Type().Elem()
	pks, err := c.getIdsByCacheSQL(sql)
	if err != nil {
		c.log.Error("getIdsByCacheSQL err: %v ,sql :%v", err, sql)
	}
	var isNoCacheSQL = len(pks) == 0
	keyValues, err := c.cacheHandler.GetAll(pks)
	if err != nil {
		c.log.Error("GetAll err: %v ,sql :%v", err, pks)
	}

	restPk := make(schemas.PK, 0)
	for _, kv := range keyValues {
		if !kv.Has {
			restPk = append(restPk, kv.Key)
			continue
		}
		entry := reflect.New(entryElementType).Interface()
		err = c.serializer.Deserialize(kv.Value, entry)
		if err != nil {
			restPk = append(restPk, kv.Key)
			continue
		}
		entriesValue = reflect.Append(entriesValue, reflect.Indirect(reflect.ValueOf(entry)))
	}

	if len(restPk) > 0 || isNoCacheSQL {
		slice := reflect.MakeSlice(entriesValue.Type(), 0, 0)
		value := reflect.New(slice.Type())
		value.Elem().Set(slice)
		err = c.databaseHandler.GetEntries(value.Interface(), sql)
		if err != nil {
			c.log.Error("GetEntries err:%v ,sql:%v", err, sql)
			return err
		}

		value = c.sort(value, pks)

		emptySlice := value.Interface()

		var res interface{}
		if gdreflect.IsPointerElement(entriesValue.Interface()) && !gdreflect.IsPointerElement(emptySlice) {
			res = gdreflect.CovertSliceStructValue2PointerValue(emptySlice)
		} else if !gdreflect.IsPointerElement(entriesValue.Interface()) && gdreflect.IsPointerElement(emptySlice) {
			res = gdreflect.CovertSlicePointerValue2StructValue(emptySlice)
		} else {
			res = emptySlice
		}

		resValue := reflect.Indirect(reflect.ValueOf(res))
		if res != nil && resValue.Len() > 0 {
			entriesValue = reflect.AppendSlice(entriesValue, resValue)
			c.storeCache(entriesValue.Interface())
		}
	}

	if isNoCacheSQL {
		err = c.setIdsByCacheSQL(restPk, sql)
		if err != nil {
			c.log.Error("setIdsByCacheSQL err:%v , restPk:%v ,sql:%v", err, restPk, sql)
		}
	}

	reflect.Indirect(reflect.ValueOf(entrySlice)).Set(entriesValue)

	return nil
}

func (c CacheHandler) GetEntry(entry interface{}) (bool, error) {
	entryParams, entryKey, err := schemas.GetEntryKey(entry.(schemas.IEntry))
	if err != nil {
		return false, err
	}

	entryValue, has, err := c.cacheHandler.Get(entryKey)
	if err != nil {
		c.log.Error("Failed to get data from cache err:%v entryKey:%v", err.Error(), entryKey)
	}
	if has {
		err = c.serializer.Deserialize(entryValue, entry)
	}

	if !has {
		has, err = c.databaseHandler.GetEntry(entry, fmt.Sprintf(builder.GetEntryByIdSQL(entry.(schemas.IEntry), entryParams)))
		if has {
			sliceValue := reflect.MakeSlice(reflect.SliceOf(reflect.Indirect(reflect.ValueOf(entry)).Type()), 0, 0)
			sliceValue = reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(entry)))
			c.storeCache(sliceValue.Interface())
		}

	}
	return has, err
}

func (c CacheHandler) DelEntries(entrySlice interface{}, sql string) error {
	err := c.GetEntries(entrySlice, sql)
	if err != nil {
		return err
	}
	pk, err := schemas.GetPKsByEntries(entrySlice)
	if err != nil {
		return err
	}
	return c.cacheHandler.DeleteAll(pk)
}

type EntryCache struct {
	entry    interface{}
	entryKey string
}

func (c CacheHandler) storeCache(entries interface{}) {
	entryCaches := make([]EntryCache, 0)
	entriesValue := reflect.Indirect(reflect.ValueOf(entries))
	for i := 0; i < entriesValue.Len(); i++ {
		_, entryKey, err := schemas.GetEntryKey(entriesValue.Index(i).Interface().(schemas.IEntry))
		if err != nil {
			continue
		}
		entryCaches = append(entryCaches, EntryCache{
			entry:    entriesValue.Index(i).Interface().(schemas.IEntry),
			entryKey: entryKey,
		})
	}

	keyValues := make([]KeyValue, 0)
	for _, entryCache := range entryCaches {
		value, err := c.serializer.Serialize(&entryCache.entry)
		if err != nil {
			c.log.Error("Failed serialize err:%v entry:%v", err, entryCache)
		}
		keyValues = append(keyValues, KeyValue{
			Key:   entryCache.entryKey,
			Value: value,
		})
	}
	err := c.cacheHandler.StoreAll(keyValues...)
	if err != nil {
		c.log.Error("Failed StoreAll err:%v keyValues:%v", err, keyValues)
	}
}

func (c CacheHandler) setIdsByCacheSQL(pks schemas.PK, sql string) error {
	data, err := c.serializer.Serialize(pks)
	if err != nil {
		return err
	}
	err = c.cacheHandler.StoreAll(KeyValue{
		Key:   sql,
		Value: data,
	})
	return err
}

func (c CacheHandler) getIdsByCacheSQL(sql string) (schemas.PK, error) {
	pks := make(schemas.PK, 0)
	ids, has, err := c.cacheHandler.Get(sql)
	if !has {
		return pks, nil
	}
	if err != nil {
		return pks, err
	}
	err = c.serializer.Deserialize(ids, &pks)
	return pks, err
}

func (c CacheHandler) sort(entriesValue reflect.Value, pks schemas.PK) reflect.Value {
	tempSliceValue := gdreflect.MakePointerSliceValue(entriesValue)
	set := make(map[string]reflect.Value)
	for i := 0; i < entriesValue.Len(); i++ {
		entry := entriesValue.Index(i).Interface()
		_, key, _ := schemas.GetEntryKey(entry.(schemas.IEntry))
		set[key] = entriesValue.Index(i)
	}

	for _, pk := range pks {
		if value, ok := set[pk]; ok {
			tempSliceValue = reflect.Append(tempSliceValue, value)
		}
	}
	return tempSliceValue
}
