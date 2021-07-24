package gdcache

import (
	"gdcache/schemas"
	"testing"
)

func TestEntryKeys_GetEntryKey(t *testing.T) {
	tests := []struct {
		name string
		es   schemas.EntryKeys
		want string
	}{
		{
			name: "",
			es: schemas.EntryKeys{
				{
					Name:  "Name",
					Param: "Peter",
				},
				{
					Name:  "City",
					Param: "Shanghai",
				},
			},
			want: "[Name:Peter]:[City:Shanghai]",
		},
		{
			name: "",
			es: schemas.EntryKeys{
				{
					Name:  "id",
					Param: "1",
				},
			},
			want: "[id:1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.es.GetEntryKey(); got != tt.want {
				t.Errorf("GetEntryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
