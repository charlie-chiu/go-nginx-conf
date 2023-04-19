package go_nginx_conf_test

import (
	"bytes"
	"fmt"
	conf "go-nginx-conf"
	c "go-nginx-conf/shortcut"
	"io/ioutil"
	"testing"
)

func TestGenerateL1Conf(t *testing.T) {
	testFixture := "test-fixture/application.l1.conf"
	expected, err := ioutil.ReadFile(testFixture)
	if err != nil {
		t.Fatalf("failed to read test fixture %q, %v", testFixture, err)
	}

	config := conf.Config{
		Directives: conf.Block{
			c.Upstream("lea_@_www_application_com_443",
				conf.SimpleDirective{Name: "server", Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}},
				conf.SimpleDirective{Name: "server", Params: c.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"}},
			),
			c.Server(
				c.Listen443SSLHTTP2,
				conf.SimpleDirective{Name: "server_name", Params: c.P{"www.application.com"}},
				conf.SimpleDirective{Name: "ssl_stapling", Params: c.P{"on"}},
				conf.SimpleDirective{Name: "ssl_stapling_verify", Params: c.P{"on"}},
				conf.SimpleDirective{Name: "ssl_certificate", Params: c.P{"/path/to/cert"}},
				conf.SimpleDirective{Name: "ssl_certificate_key", Params: c.P{"/path/to/key"}},
				conf.SimpleDirective{Name: "ssl_protocols", Params: c.P{"TLSv1", "TLSv1.1", "TLSv1.2"}},

				conf.SimpleDirective{Name: "underscores_in_headers", Params: c.P{"on"}},

				conf.SimpleDirective{Name: "proxy_set_header", Params: c.P{"X-Real-IP", "$remote_addr"}},
				conf.SimpleDirective{Name: "proxy_set_header", Params: c.P{"X-Forwarded-For", "$remote_addr"}},
				conf.SimpleDirective{Name: "proxy_set_header", Params: c.P{"X-Client-Verify", "SUCCESS"}},
				conf.SimpleDirective{Name: "proxy_set_header", Params: c.P{"X-SSL-Subject", "$ssl_client_s_dn"}},
				conf.SimpleDirective{Name: "proxy_set_header", Params: c.P{"X-SSL-Issuer", "$ssl_client_i_dn"}},

				conf.SimpleDirective{Name: "proxy_http_version", Params: c.P{"1.1"}},
				conf.SimpleDirective{Name: "proxy_connect_timeout", Params: c.P{"10"}},

				c.If("$http_user_agent ~* \"JianKongBao Monitor\"", conf.SimpleDirective{
					Name:   "return",
					Params: c.P{"200"},
				}),

				conf.SimpleDirective{Name: "error_page", Params: c.P{"497", "=307", "https://$host:$server_port$request_uri"}},
				conf.SimpleDirective{Name: "error_page", Params: c.P{"400", "414", "406", "@requestErrEvent"}},

				c.Location(c.P{"@requestErrEvent"},
					conf.SimpleDirective{Name: "root", Params: c.P{"/var/tmp/leadns_errpage"}},
					conf.SimpleDirective{Name: "ssi", Params: c.P{"on"}},
					conf.SimpleDirective{Name: "internal", Params: nil},
					conf.SimpleDirective{Name: "try_files", Params: c.P{"$uri", "/custom_error.html", "400"}},
				),
				c.Location(c.P{"~", "/purge(/.*)"},
					conf.SimpleDirective{Name: "proxy_cache_purge", Params: c.P{"hqszone", "$host$1$is_args$args"}},
				),

				conf.SimpleDirective{Name: "set", Params: c.P{"$origin_str", "34.96.119.139"}},

				c.Location(c.P{"/"},
					conf.SimpleDirective{Name: "limit_rate", Params: c.P{"68608k"}},
					conf.SimpleDirective{Name: "proxy_cache_bypass", Params: c.P{"$nocache"}},
					conf.SimpleDirective{Name: "proxy_no_cache", Params: c.P{"$nocache"}},
					conf.SimpleDirective{Name: "proxy_cache", Params: c.P{"hqszone"}},
					conf.SimpleDirective{Name: "proxy_cache_valid", Params: c.P{"301", "302", "0s"}},
					conf.SimpleDirective{Name: "proxy_cache_key", Params: c.P{"$host$uri$is_args$args$origin_str"}},
					conf.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_application_com_443"}},
				),

				conf.SimpleDirective{Name: "access_log", Params: c.P{"/var/log/nginx/hqs_access_@_www_application_com.log", "json"}},
				conf.SimpleDirective{Name: "error_log", Params: c.P{"/var/log/nginx/hqs_error_@_www_application_com.log"}},
			),
		},
	}

	actual := conf.DumpConfig(config, conf.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateSimpleDirective(t *testing.T) {
	type testCase struct {
		name   string
		input  conf.SimpleDirective
		output []byte
	}

	testCases := []testCase{
		{
			name: "listen 80",
			input: conf.SimpleDirective{
				Name:   "listen",
				Params: c.P{"80"},
			},
			output: []byte("listen 80;"),
		},
		{
			name: "custom error page",
			input: conf.SimpleDirective{
				Name:   "error_page",
				Params: c.P{"497", "=307", "https://$host:$server_port$request_uri"},
			},
			output: []byte("error_page 497 =307 https://$host:$server_port$request_uri;"),
		},
		{
			name: "set variable",
			input: conf.SimpleDirective{
				Name:   "set",
				Params: c.P{"$origin_str", "34.96.119.139"},
			},
			output: []byte("set $origin_str 34.96.119.139;"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test simple directive: %q", tc.name), func(t *testing.T) {
			actual := conf.DumpConfig(conf.Config{
				Directives: conf.Block{
					tc.input,
				},
			}, conf.IndentedStyle)

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

	config := conf.Config{
		Directives: conf.Block{
			c.Upstream("lea_@_www_application_com_443",
				conf.SimpleDirective{
					Name:   "server",
					Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"},
				},
				conf.SimpleDirective{
					Name:   "server",
					Params: c.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"},
				},
			)},
	}
	actual := conf.DumpConfig(config, conf.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateIfDirective(t *testing.T) {
	expected := []byte(
		`if ($host = 'www.application.com') {
    return 301 https://$host$request_uri;
}`)

	config := conf.Config{
		Directives: conf.Block{
			c.If(
				"$host = 'www.application.com'",
				conf.SimpleDirective{
					Name:   "return",
					Params: c.P{"301", "https://$host$request_uri"},
				},
			),
		},
	}

	actual := conf.DumpConfig(config, conf.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateLocationDirective(t *testing.T) {
	expected := []byte(`location ~ /purge(/.*) {
    proxy_cache_purge hqszone $host$1$is_args$args;
}`)

	config := conf.Config{
		Directives: conf.Block{
			c.Location(
				c.P{"~ /purge(/.*)"},
				conf.SimpleDirective{
					Name: "proxy_cache_purge", Params: c.P{"hqszone", "$host$1$is_args$args"},
				},
			),
		},
	}

	actual := conf.DumpConfig(config, conf.IndentedStyle)

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
