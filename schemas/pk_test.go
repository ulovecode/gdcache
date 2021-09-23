package schemas

import (
	"reflect"
	"testing"
)

func TestPK_ToEntryKeys(t *testing.T) {
	tests := []struct {
		name string
		pk   PK
		want []EntryKeys
	}{
		{
			name: "",
			pk:   []string{"_#[Name:Peter]-[City:Shanghai]"},
			want: []EntryKeys{
				[]EntryKey{
					{
						Name:  "Name",
						Param: "Peter",
					},
					{
						Name:  "City",
						Param: "Shanghai",
					},
				},
			},
		}, {
			name: "",
			pk:   []string{"_schemas.MockEntry#[relateId:0]-[sourceId:0]"},
			want: []EntryKeys{
				[]EntryKey{
					{
						Name:  "relateId",
						Param: "0",
					},
					{
						Name:  "sourceId",
						Param: "0",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.ToEntryKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEntryKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
