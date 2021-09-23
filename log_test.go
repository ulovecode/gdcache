package gdcache

import "testing"

func TestDefaultLogger_Debug(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				format: "1234",
				a:      []interface{}{1, "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultLogger{}
			d.Debug(tt.args.format, tt.args.a...)
		})
	}
}

func TestDefaultLogger_Error(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				format: "1234",
				a:      []interface{}{1, "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultLogger{}
			d.Error(tt.args.format, tt.args.a...)
		})
	}
}

func TestDefaultLogger_Info(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				format: "1234",
				a:      []interface{}{1, "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultLogger{}
			d.Info(tt.args.format, tt.args.a...)
		})
	}
}

func TestDefaultLogger_Warn(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				format: "1234",
				a:      []interface{}{1, "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DefaultLogger{}
			d.Warn(tt.args.format, tt.args.a...)
		})
	}
}
