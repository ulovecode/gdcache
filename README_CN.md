<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [gdcache](#gdcache)
  - [核心原理](#%E6%A0%B8%E5%BF%83%E5%8E%9F%E7%90%86)
  - [节省内存](#%E8%8A%82%E7%9C%81%E5%86%85%E5%AD%98)
  - [例子](#%E4%BE%8B%E5%AD%90)
      - [实现 ICache 接口](#%E5%AE%9E%E7%8E%B0-icache-%E6%8E%A5%E5%8F%A3)
    - [gorm 使用](#gorm-%E4%BD%BF%E7%94%A8)
    - [xorm 使用](#xorm-%E4%BD%BF%E7%94%A8)
    - [原生SQL 使用](#%E5%8E%9F%E7%94%9Fsql-%E4%BD%BF%E7%94%A8)
  - [贡献](#%E8%B4%A1%E7%8C%AE)
  - [许可证](#%E8%AE%B8%E5%8F%AF%E8%AF%81)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# gdcache

gdcache 是一个由 golang 实现的纯非侵入式分布式缓存库，你可以用它来实现你自己的分布式缓存。 [英文文档](https://github.com/ulovecode/gdcache/blob/main/README.md)

## 核心原理

gdcache 的核心原理是将 sql 转换成 id 并缓存起来。然后查询并缓存每个 id 。这样每个 sql 就可以使用每个 id 对应的实体内容了。

![img.png](doc/flow-img.png)

如上图所示，每一段sql都可以转换为对应的sql，底层去复用这些id。如果这些这些id没有被查询到，由于我们不知道到底是因为过期了，还是因为这些值在数据库中不存在，我们都会在数据库中，将这些无法从cache中取的实体从从数据库中再访问一遍获取，如果能够获取到，会进行一次缓存。

## 节省内存

常规的缓存框架，会缓存结果的内容，但 gdcache 缓存库与之不同，他只会缓存结果的id，并通过id去寻找值。这样的好处是，可以重复利用值，id 对应的值只会被缓存一次。

## 例子

要被缓存的类必须要实现 `TableName()` 方法

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

#### 实现 ICache 接口
并且实现 `ICache` 接口，可以使用redis或者gocache作为底层实现。

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

### gorm 使用

实现 `IDB` 接口

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

### xorm 使用

实现 `IDB` 接口

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
### 原生SQL 使用

实现 `IDB` 接口

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

## 贡献

您可以帮助提供更好的 gdcahe ，通过提交 pr 的方式。

## 许可证
© Jovanzhu, 2021~time.Now

发布 [MIT License](https://github.com/ulovecode/gdcache/blob/main/License)