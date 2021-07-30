package gorm

import (
	"github.com/ulovecode/gdcache"
	"github.com/ulovecode/gdcache/schemas"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MemoryCacheHandler struct {
	data map[string][]byte
}

func (m MemoryCacheHandler) StoreAll(keyValues ...gdcache.KeyValue) (err error) {
	for _, keyValue := range keyValues {
		m.data[keyValue.Key] = keyValue.Value
	}
	return nil
}

func (m MemoryCacheHandler) Get(key string) (data []byte, has bool, err error) {
	bytes, has := m.data[key]
	return bytes, has, nil
}

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

func (m MemoryCacheHandler) DeleteAll(keys schemas.PK) error {
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

func NewMemoryCacheHandler() *MemoryCacheHandler {
	return &MemoryCacheHandler{
		data: make(map[string][]byte, 0),
	}
}

type GormDB struct {
	db *gorm.DB
}

func (g GormDB) GetEntries(entries interface{}, sql string) (interface{}, error) {
	users := make([]User, 0)
	tx := g.db.Raw(sql).Find(entries)
	return users, tx.Error
}

func (g GormDB) GetEntry(entry interface{}, sql string) (interface{}, bool, error) {
	tx := g.db.Model(&entry).Raw(sql).Find(&entry)
	return entry, true, tx.Error
}

func NewGormCacheHandler() *gdcache.CacheHandler {
	return gdcache.NewCacheHandler(NewMemoryCacheHandler(), NewGormDd())
}

func NewGormDd() gdcache.IDB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return GormDB{
		db: db,
	}
}

type User struct {
	Id   uint64 `gorm:"primary_key"`
	Name string
	Age  int
}

func (u User) TableName() string {
	return "user"
}
