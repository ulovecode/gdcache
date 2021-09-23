package schemas

import (
	"fmt"
	"github.com/ulovecode/gdcache/tag"
	"testing"
)

func init() {
	tag.ConfigTag("cache")
}

// User User Info
type User struct {
	Id   uint64 `json:"id"`
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

func TestGetEntryKey(t *testing.T) {
	entry := &MockEntry{}
	userEntry := &User{}

	type args struct {
		entry IEntry
	}
	tests := []struct {
		name    string
		args    args
		want    []EntryKey
		want1   string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				entry: entry,
			},
			want: []EntryKey{{
				Name:  "relateId",
				Param: "0",
			}, {
				Name:  "sourceId",
				Param: "0",
			}},
			want1:   "_schemas.MockEntry#[relateId:0]-[sourceId:0]",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				entry: userEntry,
			},
			want: []EntryKey{{
				Name:  "id",
				Param: "0",
			}},
			want1:   "_schemas.User#[id:0]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetEntryKey(tt.args.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !(fmt.Sprint(got) == fmt.Sprint(tt.want)) {
				t.Errorf("GetEntryKey() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetEntryKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
