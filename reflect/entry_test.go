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

func TestIsPointerElement(t *testing.T) {
	structEntries := make([]MockEntry, 0)
	pointerEntries := make([]*MockEntry, 0)
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				value: &structEntries,
			},
			want: false,
		},
		{
			name: "",
			args: args{
				value: &pointerEntries,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPointerElementSlice(tt.args.value); got != tt.want {
				t.Errorf("IsPointerElementSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSliceValue1(t *testing.T) {
	structEntries := make([]MockEntry, 0)
	pointerEntries := make([]*MockEntry, 0)
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "structEntries",
			args: args{
				value: &structEntries,
			},
		},
		{
			name: "pointerEntries",
			args: args{
				value: &pointerEntries,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSliceValue(tt.args.value)
			if tt.name == "structEntries" {
				if _, ok := got.(MockEntry); ok {
					t.Errorf("GetSliceValue() = %v", got)
				}
			}
			if tt.name == "pointerEntries" {
				if _, ok := got.(*MockEntry); ok {
					t.Errorf("GetSliceValue() = %v", got)
				}
			}
		})
	}
}

func TestMakePointerSliceValue(t *testing.T) {
	structEntries := make([]MockEntry, 0)
	pointerEntries := make([]*MockEntry, 0)
	type args struct {
		entriesValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want reflect.Value
	}{
		{
			name: "",
			args: args{
				entriesValue: reflect.ValueOf(structEntries),
			},
			want: reflect.ValueOf(&structEntries),
		},
		{
			name: "",
			args: args{
				entriesValue: reflect.ValueOf(pointerEntries),
			},
			want: reflect.ValueOf(&pointerEntries),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePointerSliceValue(tt.args.entriesValue); !(fmt.Sprint(got) == fmt.Sprint(tt.want)) {
				t.Errorf("MakePointerSliceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
