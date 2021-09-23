package xorm

import (
	"github.com/go-xorm/xorm"
	"github.com/ulovecode/gdcache"
	"github.com/ulovecode/gdcache/schemas"
)

// MemoryCacheHandler memory cache handler
type MemoryCacheHandler struct {
	data map[string][]byte
}

// StoreAll Store key value
func (m MemoryCacheHandler) StoreAll(keyValues ...gdcache.KeyValue) (err error) {
	for _, keyValue := range keyValues {
		m.data[keyValue.Key] = keyValue.Value
	}
	return nil
}

// Get Get value by key
func (m MemoryCacheHandler) Get(key string) (data []byte, has bool, err error) {
	bytes, has := m.data[key]
	return bytes, has, nil
}

// GetAll Get values by keys
func (m MemoryCacheHandler) GetAll(keys schemas.PK) (data []gdcache.ReturnKeyValue, err error) {
	returnKeyValues := make([]gdcache.ReturnKeyValue, 0)
	for _, key := range keys {
		bytes, has := m.data[key]
		returnKeyValues = append(returnKeyValues, gdcache.ReturnKeyValue{
			KeyValue: gdcache.KeyValue{
				Key:   key,
				Value: bytes,
			},
			Has: has,
		})
	}
	return returnKeyValues, nil
}

// DeleteAll Delete all key caches
func (m MemoryCacheHandler) DeleteAll(keys schemas.PK) error {
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

// NewMemoryCacheHandler Create a cache handler
func NewMemoryCacheHandler() *MemoryCacheHandler {
	return &MemoryCacheHandler{
		data: make(map[string][]byte, 0),
	}
}

// XormDB  Gorm database
type XormDB struct {
	db *xorm.Engine
}

// GetEntries Get the list of entities through sql
func (g XormDB) GetEntries(entries interface{}, sql string) error {
	err := g.db.SQL(sql).Find(entries)
	return err
}

// GetEntry Get entities through sql
func (g XormDB) GetEntry(entry interface{}, sql string) (bool, error) {
	has, err := g.db.SQL(sql).Get(entry)
	return has, err
}

// NewXormCacheHandler Create a gorm cache
func NewXormCacheHandler() *gdcache.CacheHandler {
	return gdcache.NewCacheHandler(NewMemoryCacheHandler(), NewXormDd())
}

// NewXormDd Create a gorm db
func NewXormDd() gdcache.IDB {
	db, err := xorm.NewEngine("mysql", "root:root@/test?charset=utf8")
	if err != nil {
		panic(err)
	}
	db.ShowSQL(true)
	return XormDB{
		db: db,
	}
}

// User User Info
type User struct {
	Id   uint64 `cache:"id" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TableName Table Name
func (u User) TableName() string {
	return "user"
}

// MockEntry Mock entity
type MockEntry struct {
	RelateId   int64 `xorm:"relateId" cache:"relateId"`
	SourceId   int64 `xorm:"sourceId"  cache:"sourceId"`
	PropertyId int64 `xorm:"propertyId"  `
}

// TableName Table Name
func (m MockEntry) TableName() string {
	return "public_relation"
}
