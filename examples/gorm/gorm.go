package gorm

import (
	"github.com/ulovecode/gdcache"
	"github.com/ulovecode/gdcache/schemas"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

// GormDB  Gorm database
type GormDB struct {
	db *gorm.DB
}

// GetEntries Get the list of entities through sql
func (g GormDB) GetEntries(entries interface{}, sql string) error {
	tx := g.db.Raw(sql).Find(entries)
	return tx.Error
}

// GetEntry Get entities through sql
func (g GormDB) GetEntry(entry interface{}, sql string) (bool, error) {
	tx := g.db.Raw(sql).Take(entry)
	if gorm.ErrRecordNotFound == tx.Error {
		return false, nil
	}
	return tx.Error != gorm.ErrRecordNotFound, tx.Error
}

// NewGormCacheHandler Create a gorm cache
func NewGormCacheHandler() *gdcache.CacheHandler {
	return gdcache.NewCacheHandler(NewMemoryCacheHandler(), NewGormDd())
}

// NewGormDd Create a gorm db
func NewGormDd() gdcache.IDB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return GormDB{
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
