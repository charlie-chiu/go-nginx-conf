package go_nginx_conf_test

import (
	"bytes"
	"fmt"
	go_nginx_conf "go-nginx-conf"
	c "go-nginx-conf/shortcut"
	"testing"
)

//func TestGenerateL1Conf(t *testing.T) {
//	testFixture := "test-fixture/application.l1.conf"
//	expected, err := ioutil.ReadFile(testFixture)
//	if err != nil {
//		t.Fatalf("failed to read test fixture %q, %v", testFixture, err)
//	}
//
//	config := go_nginx_conf.Config{
//		Directives: &go_nginx_conf.Block{Directives: []go_nginx_conf.DirectiveInterface{
//			c.Upstream("lea_@_www_jb1228_com_80",
//				go_nginx_conf.SimpleDirective{
//					Name:   "server",
//					Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"},
//				},
//				go_nginx_conf.SimpleDirective{
//					Name:   "server",
//					Params: c.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"},
//				},
//			)}},
//	}
//	actual := go_nginx_conf.DumpConfig(config, go_nginx_conf.IndentedStyle)
//
//	assertConfigEqual(t, expected, actual)
//}

func TestGenerateSimpleDirective(t *testing.T) {
	type testCase struct {
		name   string
		input  go_nginx_conf.SimpleDirective
		output []byte
	}

	testCases := []testCase{
		{
			name: "listen 80",
			input: go_nginx_conf.SimpleDirective{
				Name:   "listen",
				Params: c.P{"80"},
			},
			output: []byte("listen 80;"),
		},
		{
			name: "custom error page",
			input: go_nginx_conf.SimpleDirective{
				Name:   "error_page",
				Params: c.P{"497", "=307", "https://$host:$server_port$request_uri"},
			},
			output: []byte("error_page 497 =307 https://$host:$server_port$request_uri;"),
		},
		{
			name: "set variable",
			input: go_nginx_conf.SimpleDirective{
				Name:   "set",
				Params: c.P{"$origin_str", "34.96.119.139"},
			},
			output: []byte("set $origin_str 34.96.119.139;"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test simple directive: %q", tc.name), func(t *testing.T) {
			actual := go_nginx_conf.DumpConfig(go_nginx_conf.Config{
				Directives: &go_nginx_conf.Block{Directives: []go_nginx_conf.DirectiveInterface{
					tc.input,
				}},
			}, go_nginx_conf.IndentedStyle)

			assertConfigEqual(t, tc.output, actual)
		})
	}
}

func TestGenerateUpstream(t *testing.T) {
	expected := []byte(
		`upstream lea_@_www_application_com_443 {
    server 35.200.43.88:80 max_fails=1 fail_timeout=10s;
    server 34.92.95.215:80 max_fails=1 fail_timeout=10s;
    least_conn;
}`)

	config := go_nginx_conf.Config{
		Directives: &go_nginx_conf.Block{Directives: []go_nginx_conf.DirectiveInterface{
			c.Upstream("lea_@_www_application_com_443",
				go_nginx_conf.SimpleDirective{
					Name:   "server",
					Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"},
				},
				go_nginx_conf.SimpleDirective{
					Name:   "server",
					Params: c.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"},
				},
			)}},
	}
	actual := go_nginx_conf.DumpConfig(config, go_nginx_conf.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateIfDirective(t *testing.T) {
	expected := []byte(
		`if ($host = 'www.application.com') {
    return 301 https://$host$request_uri;
}`)

	config := go_nginx_conf.Config{
		Directives: &go_nginx_conf.Block{Directives: []go_nginx_conf.DirectiveInterface{
			c.If(
				"$host = 'www.application.com'",
				go_nginx_conf.SimpleDirective{
					Name:   "return",
					Params: c.P{"301", "https://$host$request_uri"},
				},
			),
		},
		},
	}

	actual := go_nginx_conf.DumpConfig(config, go_nginx_conf.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateLocationDirective(t *testing.T) {
	expected := []byte(`location ~ /purge(/.*) {
    proxy_cache_purge hqszone $host$1$is_args$args;
}`)

	config := go_nginx_conf.Config{
		Directives: &go_nginx_conf.Block{Directives: []go_nginx_conf.DirectiveInterface{
			c.Location(
				c.P{"~ /purge(/.*)"},
				go_nginx_conf.SimpleDirective{
					Name:   "proxy_cache_purge",
					Params: c.P{"hqszone", "$host$1$is_args$args"},
				},
			),
		},
		},
	}

	actual := go_nginx_conf.DumpConfig(config, go_nginx_conf.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func assertConfigEqual(t *testing.T, expected []byte, actual []byte) {
	t.Helper()
	if bytes.Compare(expected, actual) != 0 {
		t.Logf("failed to assert actual equal expected\n")
		t.Logf("expected:\n%s\n", expected)
		t.Logf("actual:  \n%s\n", actual)
		t.Fail()
	}
}
