package gdcache

import (
	"fmt"
	"gdcache/builder"
	"gdcache/db"
	"gdcache/log"
	gdreflect "gdcache/reflect"
	"gdcache/schemas"
	"gdcache/tag"
)

type returnKeyValue struct {
	keyValue
	has bool
}

type keyValue struct {
	key   interface{}
	value []byte
}

type ICache interface {
	StoreAll(keyValues ...keyValue) (err error)
	Get(key string) (data []byte, has bool, err error)
	GetAll(keys schemas.PK) (data []returnKeyValue, err error)
	DeleteAll(key schemas.PK) error
}

type ICacheHandler interface {
	// GetEntries cache the entity content obtained through sql, and return the entity of the array pointer type
	GetEntries(entry []schemas.IEntry, sql string) error
	// GetEntry get a pointer to an entity type and return the entity
	GetEntry(entry schemas.IEntry) (bool, error)
	// DelEntries delete the corresponding execution entity through sql,
	// Before the update, you can query the id to be deleted first, and then delete these through defer
	DelEntries(entry []schemas.IEntry, sql string) error
}

var _ ICacheHandler = CacheHandler{}

type CacheHandler struct {
	cacheHandler    ICache
	databaseHandler db.IDB
	serializer      Serializer
	log             log.Logger
}

func NewCacheHandler(cacheHandler ICache, databaseHandler db.IDB, options ...OptionsFunc) *CacheHandler {
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

func (c CacheHandler) GetEntries(entries []schemas.IEntry, sql string) error {

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
		if !kv.has {
			restPk = append(restPk, kv.key)
			continue
		}
		entry := gdreflect.GetSliceValue(entries)
		err = c.serializer.Deserialize(kv.value, entry)
		if err != nil {
			restPk = append(restPk, kv.key)
			continue
		}
		entries = append(entries, entry)
	}
	if len(restPk) > 0 {
		emptySlice := gdreflect.MakeEmptySlice(entries)
		err = c.databaseHandler.GetEntries(emptySlice, sql)
		if err != nil {
			c.log.Error("GetEntries err:%v ,sql:%v", err, sql)
			return err
		}
		entries = append(entries, emptySlice...)
		c.storeCache(emptySlice...)
	}

	if isNoCacheSQL {
		err = c.setIdsByCacheSQL(restPk, sql)
		if err != nil {
			c.log.Error("setIdsByCacheSQL err:%v , restPk:%v ,sql:%v", err, restPk, sql)
		}
	}

	return nil
}

func (c CacheHandler) GetEntry(entry schemas.IEntry) (bool, error) {
	entryParams, entryKey, err := schemas.GetEntryKey(entry)
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
		has, err = c.databaseHandler.GetEntry(entry, fmt.Sprintf(builder.GetEntryByIdSQL(entry, entryParams)))
		c.storeCache(entry)
	}

	return has, err
}

func (c CacheHandler) DelEntries(entries []schemas.IEntry, sql string) error {
	err := c.GetEntries(entries, sql)
	if err != nil {
		return err
	}
	pk, err := schemas.GetPKsByEntries(entries)
	if err != nil {
		return err
	}
	return c.cacheHandler.DeleteAll(pk)
}

type EntryCache struct {
	entry    schemas.IEntry
	entryKey string
}

func (c CacheHandler) storeCache(entries ...schemas.IEntry) {
	entryCaches := make([]EntryCache, 0)
	for _, entry := range entries {
		_, entryKey, err := schemas.GetEntryKey(entry)
		if err != nil {
			continue
		}
		entryCaches = append(entryCaches, EntryCache{
			entry:    entry,
			entryKey: entryKey,
		})
	}

	keyValues := make([]keyValue, 0)
	for _, entryCache := range entryCaches {
		value, err := c.serializer.Serialize(entryCache)
		if err != nil {
			c.log.Error("Failed serialize err:%v entry:%v", err, entryCache)
		}
		keyValues = append(keyValues, keyValue{
			key:   entryCache.entryKey,
			value: value,
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
	err = c.cacheHandler.StoreAll(keyValue{
		key:   sql,
		value: data,
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
