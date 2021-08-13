package builder

import (
	"github.com/ulovecode/gdcache/schemas"
	"github.com/ulovecode/gdcache/tag"
	"testing"
)

type MockEntry struct {
	RelateId   int64 `cache:"relateId"`
	SourceId   int64 `cache:"sourceId"`
	PropertyId int64 `cache:"propertyId"`
}

func (m MockEntry) TableName() string {
	return "public_relation"
}

type MockEntry2 struct {
	Id         int64
	RelateId   int64
	SourceId   int64
	PropertyId int64
}

func (m MockEntry2) TableName() string {
	return "public_relation"
}

func TestGetEntryByIdSQL(t *testing.T) {
	mockEntry1 := MockEntry{
		RelateId:   213,
		SourceId:   12,
		PropertyId: 2,
	}
	mockEntry2 := MockEntry2{
		Id:         421,
		RelateId:   213,
		SourceId:   12,
		PropertyId: 2,
	}
	tag.ConfigTag("cache")
	entryParams1, _, err := schemas.GetEntryKey(mockEntry1)
	entryParams2, _, err := schemas.GetEntryKey(mockEntry2)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		entry       schemas.IEntry
		entryParams []schemas.EntryKey
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				entry:       mockEntry1,
				entryParams: entryParams1,
			},
			want: "SELECT * FROM public_relation  WHERE  relateId = 213 AND sourceId = 12 AND propertyId = 2 ;",
		},
		{
			name: "",
			args: args{
				entry:       mockEntry2,
				entryParams: entryParams2,
			},
			want: "SELECT * FROM public_relation  WHERE  id = 421 ;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEntryByIdSQL(tt.args.entry, tt.args.entryParams); got != tt.want {
				t.Errorf("GetEntryByIdSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEntriesByIdSQL(t *testing.T) {
	mockEntry1 := MockEntry{
		RelateId:   213,
		SourceId:   12,
		PropertyId: 2,
	}

	mockEntry2 := MockEntry{
		RelateId:   2,
		SourceId:   4,
		PropertyId: 5,
	}
	tag.ConfigTag("cache")
	entryParams1, _, err := schemas.GetEntryKey(mockEntry1)
	entryParams2, _, err := schemas.GetEntryKey(mockEntry2)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		entry     schemas.IEntry
		entryKeys []schemas.EntryKeys
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				entry: mockEntry1,
				entryKeys: []schemas.EntryKeys{
					entryParams1,
					entryParams2,
				},
			},
			want: "SELECT * FROM public_relation  WHERE (  relateId = 213 AND sourceId = 12 AND propertyId = 2  ) OR (  relateId = 2 AND sourceId = 4 AND propertyId = 5  );",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEntriesByIdSQL(tt.args.entry, tt.args.entryKeys); got != tt.want {
				t.Errorf("GetEntriesByIdSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFmtSql(t *testing.T) {
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id = ?",
				args: []interface{}{1},
			},
			want: "",
		},
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id in ?",
				args: []interface{}{[]string{"1", "2"}},
			},
			want: "",
		},
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id =  ? and id = ? limit ?,?",
				args: []interface{}{1, "2", 0, 10},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSql(tt.args.sql, tt.args.args...); got != tt.want {
				t.Errorf("GenerateSql() = %v, want %v", got, tt.want)
			}
		})
	}
}
