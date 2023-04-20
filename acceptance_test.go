package go_nginx_conf_test

import (
	"bytes"
	"fmt"
	conf "go-nginx-conf"
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
			conf.Upstream("lea_@_www_application_com_443",
				conf.SD{Name: "server", Params: conf.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}},
				conf.SD{Name: "server", Params: conf.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"}},
			),
			conf.Server(
				conf.Listen443SSLHTTP2,
				conf.SD{Name: "server_name", Params: conf.P{"www.application.com"}},
				conf.SD{Name: "ssl_stapling", Params: conf.P{"on"}},
				conf.SD{Name: "ssl_stapling_verify", Params: conf.P{"on"}},
				conf.SD{Name: "ssl_certificate", Params: conf.P{"/path/to/cert"}},
				conf.SD{Name: "ssl_certificate_key", Params: conf.P{"/path/to/key"}},
				conf.SD{Name: "ssl_protocols", Params: conf.P{"TLSv1", "TLSv1.1", "TLSv1.2"}},

				conf.SD{Name: "underscores_in_headers", Params: conf.P{"on"}},

				conf.SD{Name: "proxy_set_header", Params: conf.P{"X-Real-IP", "$remote_addr"}},
				conf.SD{Name: "proxy_set_header", Params: conf.P{"X-Forwarded-For", "$remote_addr"}},
				conf.SD{Name: "proxy_set_header", Params: conf.P{"X-Client-Verify", "SUCCESS"}},
				conf.SD{Name: "proxy_set_header", Params: conf.P{"X-SSL-Subject", "$ssl_client_s_dn"}},
				conf.SD{Name: "proxy_set_header", Params: conf.P{"X-SSL-Issuer", "$ssl_client_i_dn"}},

				conf.SD{Name: "proxy_http_version", Params: conf.P{"1.1"}},
				conf.SD{Name: "proxy_connect_timeout", Params: conf.P{"10"}},

				conf.If("$http_user_agent ~* \"JianKongBao Monitor\"", conf.SD{
					Name:   "return",
					Params: conf.P{"200"},
				}),

				conf.SD{Name: "error_page", Params: conf.P{"497", "=307", "https://$host:$server_port$request_uri"}},
				conf.SD{Name: "error_page", Params: conf.P{"400", "414", "406", "@requestErrEvent"}},

				conf.Location(conf.P{"@requestErrEvent"},
					conf.SD{Name: "root", Params: conf.P{"/var/tmp/leadns_errpage"}},
					conf.SD{Name: "ssi", Params: conf.P{"on"}},
					conf.SD{Name: "internal", Params: nil},
					conf.SD{Name: "try_files", Params: conf.P{"$uri", "/custom_error.html", "400"}},
				),
				conf.Location(conf.P{"~", "/purge(/.*)"},
					conf.SD{Name: "proxy_cache_purge", Params: conf.P{"hqszone", "$host$1$is_args$args"}},
				),

				conf.SD{Name: "set", Params: conf.P{"$origin_str", "34.96.119.139"}},

				conf.Location(conf.P{"/"},
					conf.SD{Name: "limit_rate", Params: conf.P{"68608k"}},
					conf.SD{Name: "proxy_cache_bypass", Params: conf.P{"$nocache"}},
					conf.SD{Name: "proxy_no_cache", Params: conf.P{"$nocache"}},
					conf.SD{Name: "proxy_cache", Params: conf.P{"hqszone"}},
					conf.SD{Name: "proxy_cache_valid", Params: conf.P{"301", "302", "0s"}},
					conf.SD{Name: "proxy_cache_key", Params: conf.P{"$host$uri$is_args$args$origin_str"}},
					conf.SD{Name: "proxy_pass", Params: conf.P{"$scheme://lea_@_www_application_com_443"}},
				),

				conf.SD{Name: "access_log", Params: conf.P{"/var/log/nginx/hqs_access_@_www_application_com.log", "json"}},
				conf.SD{Name: "error_log", Params: conf.P{"/var/log/nginx/hqs_error_@_www_application_com.log"}},
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
				Params: conf.P{"80"},
			},
			output: []byte("listen 80;"),
		},
		{
			name: "custom error page",
			input: conf.SimpleDirective{
				Name:   "error_page",
				Params: conf.P{"497", "=307", "https://$host:$server_port$request_uri"},
			},
			output: []byte("error_page 497 =307 https://$host:$server_port$request_uri;"),
		},
		{
			name: "set variable",
			input: conf.SimpleDirective{
				Name:   "set",
				Params: conf.P{"$origin_str", "34.96.119.139"},
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
			conf.Upstream("lea_@_www_application_com_443",
				conf.SimpleDirective{
					Name:   "server",
					Params: conf.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"},
				},
				conf.SimpleDirective{
					Name:   "server",
					Params: conf.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"},
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
			conf.If(
				"$host = 'www.application.com'",
				conf.SimpleDirective{
					Name:   "return",
					Params: conf.P{"301", "https://$host$request_uri"},
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
			conf.Location(
				conf.P{"~ /purge(/.*)"},
				conf.SimpleDirective{
					Name: "proxy_cache_purge", Params: conf.P{"hqszone", "$host$1$is_args$args"},
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
