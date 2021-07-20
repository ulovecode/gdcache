package tags

import (
	"reflect"
	"testing"
)

type TagTest struct {
	Id  int64 `cacheId:"11"`
	Id2 int64
	Id3 int64 `cacheId:"2"`
	Id4 int64 `cacheId:"3"`
}

func TestTag_GetPkTagField(t1 *testing.T) {
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
			fields: fields{
				tagName: "cacheId",
			},
			args: args{
				value: tagTest,
			},
			want: []reflect.StructField{
				{
					Name:      "Id3",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cacheId:"2"`,
					Offset:    16,
					Index:     []int{2},
					Anonymous: false,
				},
				{
					Name:      "Id4",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cacheId:"3"`,
					Offset:    24,
					Index:     []int{3},
					Anonymous: false,
				},
				{
					Name:      "Id",
					PkgPath:   "",
					Type:      reflect.TypeOf(int64(0)),
					Tag:       `cacheId:"11"`,
					Offset:    0,
					Index:     []int{0},
					Anonymous: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Tag{
				tagName: tt.fields.tagName,
			}
			got := t.GetCacheTagSortFields(tt.args.value)

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetCacheTagSortFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}
