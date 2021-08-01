<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [gdcache](#gdcache)
  - [Core principle](#core-principle)
  - [Save memory](#save-memory)
  - [Examples](#examples)
      - [Implement the ICache interface](#implement-the-icache-interface)
    - [Gorm use](#gorm-use)
    - [Xorm use](#xorm-use)
    - [Native SQL use](#native-sql-use)
  - [Contributing](#contributing)
  - [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# gdcache

gdcache is a pure non-intrusive distributed cache library implemented by golang, you can use it to implement your own
distributed cache. [中文文档](https://github.com/ulovecode/gdcache/blob/main/README_CN.md)

## Core principle

The core principle of gdcache is to convert sql into id and cache them. Then each id is queried and cached. In this way, each sql can use the entity content corresponding to each id.

![img.png](doc/flow-img.png)

As shown in the figure above, each piece of sql can be converted to the corresponding sql, and the bottom layer will reuse these ids. If these ids are not queried, because we don’t know whether they are out of date or because these values do not exist in the database, we will all be in the database and access these entities that cannot be retrieved from the cache from the database. Get it once, and if it can get it, it will cache it once.

## Save memory

The conventional caching framework will cache the content of the result, but the gdcache cache library is different from it. It will only cache the id of the result and find the value through the id. The advantage of this is that the value can be reused, and the value corresponding to id will only be cached once.


## Examples

The class to be cached must implement the `TableName()` method

```go
type User struct {
	Id   uint64
	Name string
	Age  int
}

func (u User) TableName() string {
	return "user"
}
```

#### Implement the ICache interface

And to implement the `ICache` interface, you can use redis or gocache as the underlying implementation.

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

### Gorm use

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

### Xorm use

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
### Native SQL use

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

## Contributing

You can help provide better _gdcahe_ by submitting pr.

## License
© Jovanzhu, 2021~time.Now

Released under the  [MIT License](https://github.com/ulovecode/gdcache/blob/main/License)