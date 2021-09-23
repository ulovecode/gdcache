package gdcache

import (
	"encoding/json"
	"errors"
	"github.com/ulovecode/gdcache/schemas"
	"testing"
)

// MemoryCacheHandler memory cache handler
type MemoryCacheHandler struct {
	data map[string][]byte
}

// StoreAll Store key value
func (m MemoryCacheHandler) StoreAll(keyValues ...KeyValue) (err error) {
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
func (m MemoryCacheHandler) GetAll(keys schemas.PK) (data []ReturnKeyValue, err error) {
	returnKeyValues := make([]ReturnKeyValue, 0)
	for _, key := range keys {
		bytes, has := m.data[key]
		returnKeyValues = append(returnKeyValues, ReturnKeyValue{
			KeyValue: KeyValue{
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

var data = make(map[string][]byte, 0)

// NewMemoryCacheHandler Create a cache handler
func NewMemoryCacheHandler() *MemoryCacheHandler {
	return &MemoryCacheHandler{
		data: data,
	}
}

// MemoryDb Memory Database
type MemoryDb struct {
}

// NewMemoryDb Create Memory Database
func NewMemoryDb() *MemoryDb {
	return &MemoryDb{}
}

// GetEntries Get the list of entities through sql
func (m MemoryDb) GetEntries(entries interface{}, sql string) error {
	if sql == "SELECT * FROM public_relation  WHERE  relateId = 1 AND sourceId = 2 AND propertyId = 3 ;" {
		mockEntries := make([]MockEntry, 0)
		mockEntries = append(mockEntries, MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		})
		marshal, _ := json.Marshal(mockEntries)
		json.Unmarshal(marshal, entries)
		return nil
	} else if sql == "SELECT * FROM public_relation  WHERE  relateId = 1 AND sourceId = 2;" {
		mockEntries := make([]*MockEntry, 0)
		mockEntries = append(mockEntries, &MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		})
		mockEntries = append(mockEntries, &MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 4,
		})
		marshal, _ := json.Marshal(mockEntries)
		json.Unmarshal(marshal, entries)
		return nil
	} else if sql == "SELECT * FROM public_relation  WHERE (  relateId = 1 AND sourceId = 2 AND propertyId = 3  ) OR (  relateId = 1 AND sourceId = 2 AND propertyId = 4  );" {
		mockEntries := make([]MockEntry, 0)
		mockEntries = append(mockEntries, MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		})
		mockEntries = append(mockEntries, MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 4,
		})
		marshal, _ := json.Marshal(mockEntries)
		json.Unmarshal(marshal, entries)
		return nil
	}
	return errors.New("mockEntries not found")
}

// GetEntry Get entities through sql
func (m MemoryDb) GetEntry(entry interface{}, sql string) (bool, error) {
	if sql == "SELECT * FROM public_relation  WHERE  relateId = 1 AND sourceId = 2 AND propertyId = 3 ;" {
		mockEntry := &MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		}
		marshal, _ := json.Marshal(mockEntry)
		json.Unmarshal(marshal, entry)
		return true, nil
	} else if sql == "SELECT COUNT(*) FROM (SELECT * FROM public_relation  WHERE  relateId = 1 AND sourceId = 2 AND propertyId = 3 ;) t" {
		i := entry.(*int64)
		*i = 1
		return true, nil
	}
	return false, nil
}

// NewMemoryCache Create a memory cache
func NewMemoryCache() *CacheHandler {
	return NewCacheHandler(NewMemoryCacheHandler(), NewMemoryDb())
}

// MockEntry Mock entity
type MockEntry struct {
	RelateId   int64 `cache:"relateId"`
	SourceId   int64 `cache:"sourceId"`
	PropertyId int64 `cache:"propertyId"`
}

// TableName Table Name
func (m MockEntry) TableName() string {
	return "public_relation"
}

func TestCacheHandler_GetEntry(t *testing.T) {
	type fields struct {
		cacheHandler    ICache
		databaseHandler IDB
		serializer      Serializer
		log             Logger
	}
	type args struct {
		entry interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entry: MockEntry{
					RelateId:   1,
					SourceId:   2,
					PropertyId: 3,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CacheHandler{
				cacheHandler:    tt.fields.cacheHandler,
				databaseHandler: tt.fields.databaseHandler,
				serializer:      tt.fields.serializer,
				log:             tt.fields.log,
			}
			got, err := c.GetEntry(tt.args.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEntry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheHandler_GetEntries(t *testing.T) {
	mockEntries := make([]MockEntry, 0)
	type fields struct {
		cacheHandler    ICache
		databaseHandler IDB
		serializer      Serializer
		log             Logger
	}
	type args struct {
		entrySlice interface{}
		sql        string
		args       []interface{}
	}
	memoryCacheHandler := NewMemoryCacheHandler()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				cacheHandler:    memoryCacheHandler,
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entrySlice: &mockEntries,
				sql:        "SELECT * FROM public_relation  WHERE  relateId = ? AND sourceId = ? AND propertyId = ? ;",
				args:       []interface{}{1, 2, 3},
			},
			wantErr: false,
		},
		{
			name: "",
			fields: fields{
				cacheHandler:    memoryCacheHandler,
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entrySlice: &mockEntries,
				sql:        "SELECT * FROM public_relation  WHERE  relateId = ? AND sourceId = ? AND propertyId = ? ;",
				args:       []interface{}{1, 2, 3},
			},
			wantErr: false,
		},
		{
			name: "",
			fields: fields{
				cacheHandler:    memoryCacheHandler,
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entrySlice: &mockEntries,
				sql:        "SELECT * FROM public_relation  WHERE  relateId = ? AND sourceId = ?;",
				args:       []interface{}{1, 2},
			},
			wantErr: false,
		}, {
			name: "restPk",
			fields: fields{
				cacheHandler:    memoryCacheHandler,
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entrySlice: &mockEntries,
				sql:        "SELECT * FROM public_relation  WHERE  relateId = ? AND sourceId = ?;",
				args:       []interface{}{1, 2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CacheHandler{
				cacheHandler:    tt.fields.cacheHandler,
				databaseHandler: tt.fields.databaseHandler,
				serializer:      tt.fields.serializer,
				log:             tt.fields.log,
			}
			if tt.name == "restPk" {
				println(data)
				delete(data, "_gdcache.MockEntry#[relateId:1]-[sourceId:2]-[propertyId:4]")
			}
			if err := c.GetEntries(tt.args.entrySlice, tt.args.sql, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("GetEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCacheHandler_GetEntriesAndCount(t *testing.T) {
	mockEntries := make([]MockEntry, 0)
	type fields struct {
		cacheHandler    ICache
		databaseHandler IDB
		serializer      Serializer
		log             Logger
	}
	type args struct {
		entries interface{}
		sql     string
		args    []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			fields: fields{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entries: &mockEntries,
				sql:     "SELECT * FROM public_relation  WHERE  relateId = ? AND sourceId = ? AND propertyId = ? ;",
				args:    []interface{}{1, 2, 3},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CacheHandler{
				cacheHandler:    tt.fields.cacheHandler,
				databaseHandler: tt.fields.databaseHandler,
				serializer:      tt.fields.serializer,
				log:             tt.fields.log,
			}
			got, err := c.GetEntriesAndCount(tt.args.entries, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntriesAndCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEntriesAndCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheHandler_DelEntries(t *testing.T) {
	mockEntries := make([]MockEntry, 0)
	type fields struct {
		cacheHandler    ICache
		databaseHandler IDB
		serializer      Serializer
		log             Logger
	}
	type args struct {
		entrySlice interface{}
		sql        string
		args       []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entrySlice: &mockEntries,
				sql:        "SELECT * FROM public_relation  WHERE  relateId = ? AND sourceId = ? AND propertyId = ? ;",
				args:       []interface{}{1, 2, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CacheHandler{
				cacheHandler:    tt.fields.cacheHandler,
				databaseHandler: tt.fields.databaseHandler,
				serializer:      tt.fields.serializer,
				log:             tt.fields.log,
			}
			if err := c.DelEntries(tt.args.entrySlice, tt.args.sql, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("DelEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type TestLogger struct {
}

func (t TestLogger) Info(format string, a ...interface{}) {
	panic("implement me")
}

func (t TestLogger) Error(format string, a ...interface{}) {
	panic("implement me")
}

func (t TestLogger) Debug(format string, a ...interface{}) {
	panic("implement me")
}

func (t TestLogger) Warn(format string, a ...interface{}) {
	panic("implement me")
}

type TestSerializer struct {
}

func (t TestSerializer) Serialize(value interface{}) ([]byte, error) {
	panic("implement me")
}

func (t TestSerializer) Deserialize(data []byte, ptr interface{}) error {
	panic("implement me")
}

func TestNewCacheHandler(t *testing.T) {
	type args struct {
		cacheHandler    ICache
		databaseHandler IDB
		options         []OptionsFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				options:         []OptionsFunc{WithServiceName("test")},
			},
		},
		{
			name: "",
			args: args{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				options:         []OptionsFunc{WithServiceName("test"), WithSerializer(TestSerializer{}), WithCacheTagName("test"), WithLogger(TestLogger{})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCacheHandler(tt.args.cacheHandler, tt.args.databaseHandler, tt.args.options...); !(schemas.ServiceName == "test") {
				t.Errorf("NewCacheHandler() = %v", got)
			}
		})
	}
}
