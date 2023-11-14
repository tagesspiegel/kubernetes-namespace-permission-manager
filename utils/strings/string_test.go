package strings

import (
	"reflect"
	"testing"
)

func TestKeyValue(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "valid key value string",
			args:      args{s: "foo=bar"},
			wantKey:   "foo",
			wantValue: "bar",
			wantErr:   false,
		},
		{
			name:      "invalid key value string",
			args:      args{s: "foo"},
			wantKey:   "",
			wantValue: "",
			wantErr:   true,
		},
		{
			name:      "invalid key value string with multiple =",
			args:      args{s: "foo=bar=baz"},
			wantKey:   "",
			wantValue: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, err := KeyValue(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("KeyValue() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("KeyValue() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestProperties(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid properties string",
			args: args{s: "foo=bar;bar=baz"},
			want: []string{"foo=bar", "bar=baz"},
		},
		{
			name: "invalid properties string",
			args: args{s: "foo=bar;bar=baz;"},
			want: []string{"foo=bar", "bar=baz", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Properties(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Properties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArray(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid array string",
			args: args{s: "foo,bar,baz"},
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "valid array string with empty string",
			args: args{s: "foo,bar,baz,"},
			want: []string{"foo", "bar", "baz", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Array(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveEmpty(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid array with empty string",
			args: args{s: []string{"foo", "", "bar"}},
			want: []string{"foo", "bar"},
		},
		{
			name: "valid array without empty string",
			args: args{s: []string{"foo", "bar"}},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveEmpty(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayC(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid array string",
			args: args{s: "foo::bar::baz"},
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "valid array string with empty string",
			args: args{s: "foo::bar::baz::"},
			want: []string{"foo", "bar", "baz", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayC(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayC() = %v, want %v", got, tt.want)
			}
		})
	}
}
