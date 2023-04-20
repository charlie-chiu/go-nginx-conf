package go_nginx_conf

import (
	"reflect"
	"testing"
)

func TestLocation(t *testing.T) {
	type args struct {
		parameters []string
		directives []Directive
	}
	tests := []struct {
		name string
		args args
		want BlockDirective
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Location(tt.args.parameters, tt.args.directives...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Location() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer(t *testing.T) {
	type args struct {
		directives []Directive
	}
	tests := []struct {
		name string
		args args
		want BlockDirective
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Server(tt.args.directives...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpstream(t *testing.T) {
	t.Parallel()
	type args struct {
		upstream string
		servers  []SimpleDirective
	}
	tests := []struct {
		name string
		args args
		want BlockDirective
	}{
		{
			name: "test attach `least_conn` directive",
			args: args{
				upstream: "dummy_upstream",
				servers:  nil,
			},
			want: BlockDirective{
				Name:   "upstream",
				Params: []string{"dummy_upstream"},
				Block: []Directive{
					SD{Name: "least_conn"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Upstream(tt.args.upstream, tt.args.servers...); !reflect.DeepEqual(got, tt.want) {
				t.Logf("got: \n%s", DumpDirective(got, NoIndentStyle))
				t.Logf("want:\n%s", DumpDirective(tt.want, NoIndentStyle))
				t.Fail()
			}
		})
	}
}
