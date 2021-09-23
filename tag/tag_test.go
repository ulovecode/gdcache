package tag

import (
	"reflect"
	"testing"
)

type TagTest struct {
	Id  int64 `cache:"11"`
	Id2 int64
	Id3 int64 `cache:"2"`
	Id4 int64 `cache:"3"`
	Id5 int64 `cache:"1"`
}

func TestTag_GetPkTagField(t1 *testing.T) {
	ConfigTag("cache")
	type fields struct {
		tagName string
	}
	type args struct {
		value interface{}
	}
	tagTest := TagTest{
		Id: 1,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []reflect.StructField
		wantErr bool
	}{
		{
			name: "",
			args: args{
				value: tagTest,
			},
			want: []reflect.StructField{
				{
					Name:      "Id",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cache:"11"`,
					Offset:    0,
					Index:     []int{0},
					Anonymous: false,
				},
				{
					Name:      "Id3",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cache:"2"`,
					Offset:    16,
					Index:     []int{2},
					Anonymous: false,
				},
				{
					Name:      "Id4",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cache:"3"`,
					Offset:    24,
					Index:     []int{3},
					Anonymous: false,
				},
				{
					Name:      "Id5",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cache:"1"`,
					Offset:    32,
					Index:     []int{4},
					Anonymous: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got := GetCacheTagFields(tt.args.value)

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetCacheTagFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "cache",
			want: "cache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
