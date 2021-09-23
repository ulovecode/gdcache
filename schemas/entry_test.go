package schemas

import (
	"reflect"
	"testing"
)

func TestGetPKsByEntries(t *testing.T) {
	mockEntries := make([]MockEntry, 0)
	mockEntries = append(mockEntries, MockEntry{
		RelateId:   1,
		SourceId:   2,
		PropertyId: 3,
	})
	type args struct {
		entries interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    PK
		wantErr bool
	}{
		{
			name: "",
			args: args{
				entries: mockEntries,
			},
			want:    []string{"_schemas.MockEntry#[relateId:1]-[sourceId:2]"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPKsByEntries(tt.args.entries)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPKsByEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPKsByEntries() got = %v, want %v", got, tt.want)
			}
		})
	}
}
