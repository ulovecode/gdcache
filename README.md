<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [gdcache](#gdcache)
  - [Features](#features)
  - [Core principle](#core-principle)
  - [Save memory](#save-memory)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [Gorm usage](#gorm-usage)
    - [Xorm usage](#xorm-usage)
    - [Native SQL usage](#native-sql-usage)
  - [How to use](#how-to-use)
  - [Contributing](#contributing)
  - [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# gdcache

gdcache is a pure non-intrusive cache library implemented by golang, you can use it to implement your own
cache. [中文文档](https://github.com/ulovecode/gdcache/blob/main/README_CN.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/ulovecode/gdcache)](https://goreportcard.com/report/github.com/ulovecode/gdcache)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ulovecode/gdcache)](https://pkg.go.dev/github.com/ulovecode/gdcache)
[![codecov](https://codecov.io/gh/ulovecode/gdcache/branch/main/graph/badge.svg?token=4GNQINA6QX)](https://codecov.io/gh/ulovecode/gdcache)

## Features

- Automatically cache sql
- Reuse id cache
- Adapt to Xorm and Gorm framework
- Support cache joint key
- Lightweight
- Non-invasive
- High performance
- Flexible

## Core principle

The core principle of gdcache is to convert sql into id and cache it, and cache the entity corresponding to id. In this way, each sql has the same id and can reuse the corresponding entity content.

![img.png](doc/flow-img.png)

As shown in the figure above, each piece of sql can be converted to the corresponding sql, and the bottom layer will reuse these ids. If these ids are not queried, because we don’t know whether they are out of date or because these values do not exist in the database, we will all be in the database and access these entities that cannot be retrieved from the cache from the database. Get it once, and if it can get it, it will cache it once.

## Save memory

The conventional caching framework will cache the content of the result, but the gdcache cache library is different from it. It will only cache the id of the result and find the value through the id. The advantage of this is that the value can be reused, and the value corresponding to id will only be cached once.

## Installation

```shell
go get github.com/ulovecode/gdcache
```

## Quick Start

- The class to be cached must implement the `TableName()` method, and use `cache:"id"` to indicate the cached key. The default is to cache by `id`, and the value of the `cache` tag corresponds to Fields in the database, usually can be ignored.

```go
type User struct {
	Id   uint64 `cache:"id"` // Or omit the tag
	Name string 
	Age  int
}

func (u User) TableName() string {
	return "user"
}
```

- If you want to use a joint key, you can add a `cache` tag to multiple fields

```go
type PublicRelations struct {
      RelatedId uint64 `cache:"related_id"`
      RelatedType string
      SourceId uint64 `cache:"source_id"`
      SourceType string
}

func (u PublicRelations) TableName() string {
    return "public_relations"
}
```

- Implement the `ICache` interface, you can use redis or gocache as the underlying implementation.

```go
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
```

### Gorm usage

Implement the `IDB` interface

```go
type GormDB struct {
	db *gorm.DB
}

func (g GormDB) GetEntries(entries interface{}, sql string) error {
	tx := g.db.Raw(sql).Find(entries)
	return tx.Error
}

func (g GormDB) GetEntry(entry interface{}, sql string) (bool, error) {
    tx := g.db.Raw(sql).Take(entry)
    if gorm.ErrRecordNotFound == tx.Error {
    	return false, nil
    }
    return tx.Error != gorm.ErrRecordNotFound, tx.Error
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
```

### Xorm usage

Implement the `IDB` interface

```go
type XormDB struct {
	db *xorm.Engine
}

func (g XormDB) GetEntries(entries interface{}, sql string) ( error) {
	err := g.db.SQL(sql).Find(entries)
	return  err
}

func (g XormDB) GetEntry(entry interface{}, sql string) ( bool, error) {
	has, err := g.db.SQL(sql).Get(entry)
	return has, err
}

func NewXormCacheHandler() *gdcache.CacheHandler {
	return gdcache.NewCacheHandler(NewMemoryCacheHandler(), NewXormDd())
}

func NewXormDd() gdcache.IDB {
	db, err := xorm.NewEngine("mysql", "root:root@/test?charset=utf8")
	if err != nil {
		panic(err)
	}
	return XormDB{
		db: db,
	}
}
```
### Native SQL usage

Implement the `IDB` interface

```go
type MemoryDb struct {
}

func NewMemoryDb() *MemoryDb {
	return &MemoryDb{}
}

func (m MemoryDb) GetEntries(entries interface{}, sql string) error {
	mockEntries := make([]MockEntry, 0)
	mockEntries = append(mockEntries, MockEntry{
		RelateId:   1,
		SourceId:   2,
		PropertyId: 3,
	})
	marshal, _ := json.Marshal(mockEntries)
	json.Unmarshal(marshal, entries)
	return nil
}

func (m MemoryDb) GetEntry(entry interface{}, sql string) (bool, error) {
	mockEntry := &MockEntry{
		RelateId:   1,
		SourceId:   2,
		PropertyId: 3,
	}
	marshal, _ := json.Marshal(mockEntry)
	json.Unmarshal(marshal, entry)
	return true, nil
}

func NewMemoryCache() *gdcache.CacheHandler {
	return gdcache.NewCacheHandler(NewMemoryCacheHandler(), NewMemoryDb())
}
```


## How to use

When querying a single entity, query through the entity's id and fill it into the entity. When getting multiple entities, you can use any sql query and finally fill it into the entity. Both methods must be introduced into the body's pointer.
```go
func TestNewGormCache(t *testing.T) {

	handler := NewGormCacheHandler()

	user := User{
		Id: 1,
	}
	has, err := handler.GetEntry(&user)
	if err != nil {
		t.FailNow()
	}
	if has {
		t.Logf("%v", user)
	}

	users := make([]User, 0)
	err = handler.GetEntries(&users, "SELECT * FROM user WHERE name = '33'")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users {
		t.Logf("%v", user)
	}

	err = handler.GetEntries(&users, "SELECT * FROM user WHERE id in (3)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users {
		t.Logf("%v", user)
	}
	
        count, err = handler.GetEntriesAndCount(&users1, "SELECT * FROM user WHERE id in (1,2)")
        if err != nil {
           t.FailNow()
        }
        for _, user := range users1 {
          t.Logf("%v", user)
        }
        t.Log(count)
}
        users3 := make([]User, 0)
        ids := make([]uint64, 0)
        count, err = handler.GetEntriesAndCount(&users3, "SELECT * FROM user WHERE id in ?", ids)
        if err != nil {
           t.FailNow()
        }
        for _, user := range users1 {
           t.Logf("%v", user)
        }
        t.Log(count)
        
        
        count, err = handler.GetEntriesAndCount(&users1, "SELECT * FROM user WHERE id =  ?", 1)
        if err != nil {
        	t.FailNow()
        }
        for _, user := range users1 {
        	t.Logf("%v", user)
        }
        t.Log(count)
```

Support placeholder `?`, replacement arrays and basic types

## Contributing

You can help provide better _gdcahe_ by submitting pr.

## License
© Jovanzhu, 2021~time.Now

Released under the  [MIT License](https://github.com/ulovecode/gdcache/blob/main/LICENSE)
