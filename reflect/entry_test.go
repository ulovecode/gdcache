package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

type MockEntry struct {
	RelateId   int64 `cache:"relateId"`
	SourceId   int64 `cache:"sourceId"`
	PropertyId int64 `cache:"propertyId"`
}

func (m MockEntry) GetTableName() string {
	return "public_relation"
}

func TestGetSliceValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "",
			args: args{
				value: &[]MockEntry{
					{
						RelateId:   1,
						SourceId:   2,
						PropertyId: 3,
					},
				},
			},
			want: &MockEntry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSliceValue(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSliceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCovertSlicePointerValue2StructValue(t *testing.T) {
	entries := make([]*MockEntry, 0)
	entries = append(entries, &MockEntry{
		RelateId:   1,
		SourceId:   2,
		PropertyId: 3,
	})
	entries = append(entries, &MockEntry{
		RelateId:   2,
		SourceId:   3,
		PropertyId: 4,
	})
	mockEntries := CovertSlicePointerValue2StructValue(entries).([]MockEntry)
	for i := range mockEntries {
		fmt.Printf("%v", mockEntries[i])
	}
}

func TestCovertSliceStructValue2PointerValue(t *testing.T) {
	entries := make([]MockEntry, 0)
	entries = append(entries, MockEntry{
		RelateId:   1,
		SourceId:   2,
		PropertyId: 3,
	})
	entries = append(entries, MockEntry{
		RelateId:   2,
		SourceId:   3,
		PropertyId: 4,
	})
	mockEntries := CovertSliceStructValue2PointerValue(entries).([]*MockEntry)
	for i := range mockEntries {
		fmt.Printf("%v", mockEntries[i])
	}
}
