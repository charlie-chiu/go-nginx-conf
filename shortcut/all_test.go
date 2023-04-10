package shortcut

import (
	go_nginx_conf "go-nginx-conf"
	"reflect"
	"testing"
)

func TestLocation(t *testing.T) {
	type args struct {
		parameters []string
		directives []go_nginx_conf.Directive
	}
	tests := []struct {
		name string
		args args
		want go_nginx_conf.BlockDirective
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
		directives []go_nginx_conf.Directive
	}
	tests := []struct {
		name string
		args args
		want go_nginx_conf.BlockDirective
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
		servers  []go_nginx_conf.SimpleDirective
	}
	tests := []struct {
		name string
		args args
		want go_nginx_conf.BlockDirective
	}{
		{
			name: "test attach `least_conn` directive",
			args: args{
				upstream: "dummy_upstream",
				servers:  nil,
			},
			want: go_nginx_conf.BlockDirective{
				Name:   "upstream",
				Params: []string{"dummy_upstream"},
				Block: &go_nginx_conf.Block{Directives: []go_nginx_conf.Directive{
					go_nginx_conf.SimpleDirective{Name: "least_conn"},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Upstream(tt.args.upstream, tt.args.servers...); !reflect.DeepEqual(got, tt.want) {
				t.Logf("got: \n%s", go_nginx_conf.DumpDirective(got, go_nginx_conf.NoIndentStyle))
				t.Logf("want:\n%s", go_nginx_conf.DumpDirective(tt.want, go_nginx_conf.NoIndentStyle))
				t.Fail()
			}
		})
	}
}
