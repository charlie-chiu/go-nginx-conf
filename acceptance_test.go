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
				conf.SD("server", conf.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}, nil),
				conf.SD("server", conf.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"}, nil),
			),
			conf.Server(
				conf.Listen443SSLHTTP2,
				conf.SD("server_name", conf.P{"www.application.com"}, nil),
				conf.SD("ssl_stapling", conf.P{"on"}, nil),
				conf.SD("ssl_stapling_verify", conf.P{"on"}, nil),
				conf.SD("ssl_certificate", conf.P{"/path/to/cert"}, nil),
				conf.SD("ssl_certificate_key", conf.P{"/path/to/key"}, nil),
				conf.SD("ssl_protocols", conf.P{"TLSv1", "TLSv1.1", "TLSv1.2"}, nil),

				conf.SD("underscores_in_headers", conf.P{"on"}, nil),

				conf.SD("proxy_set_header", conf.P{"X-Real-IP", "$remote_addr"}, nil),
				conf.SD("proxy_set_header", conf.P{"X-Forwarded-For", "$remote_addr"}, nil),
				conf.SD("proxy_set_header", conf.P{"X-Client-Verify", "SUCCESS"}, nil),
				conf.SD("proxy_set_header", conf.P{"X-SSL-Subject", "$ssl_client_s_dn"}, nil),
				conf.SD("proxy_set_header", conf.P{"X-SSL-Issuer", "$ssl_client_i_dn"}, nil),

				conf.SD("proxy_http_version", conf.P{"1.1"}, nil),
				conf.SD("proxy_connect_timeout", conf.P{"10"}, nil),

				conf.If("$http_user_agent ~* \"JianKongBao Monitor\"", conf.SimpleDirective{
					Name:   "return",
					Params: conf.P{"200"},
				}),

				conf.SD("error_page", conf.P{"497", "=307", "https://$host:$server_port$request_uri"}, nil),
				conf.SD("error_page", conf.P{"400", "414", "406", "@requestErrEvent"}, nil),

				conf.Location(conf.P{"@requestErrEvent"},
					conf.SD("root", conf.P{"/var/tmp/leadns_errpage"}, nil),
					conf.SD("ssi", conf.P{"on"}, nil),
					conf.SD("internal", nil, nil),
					conf.SD("try_files", conf.P{"$uri", "/custom_error.html", "400"}, nil),
				),
				conf.Location(conf.P{"~", "/purge(/.*)"},
					conf.SD("proxy_cache_purge", conf.P{"hqszone", "$host$1$is_args$args"}, nil),
				),

				conf.SD("set", conf.P{"$origin_str", "34.96.119.139"}, nil),

				conf.Location(conf.P{"/"},
					conf.SD("limit_rate", conf.P{"68608k"}, nil),
					conf.SD("proxy_cache_bypass", conf.P{"$nocache"}, nil),
					conf.SD("proxy_no_cache", conf.P{"$nocache"}, nil),
					conf.SD("proxy_cache", conf.P{"hqszone"}, nil),
					conf.SD("proxy_cache_valid", conf.P{"301", "302", "0s"}, nil),
					conf.SD("proxy_cache_key", conf.P{"$host$uri$is_args$args$origin_str"}, nil),
					conf.SD("proxy_pass", conf.P{"$scheme://lea_@_www_application_com_443"}, nil),
				),

				conf.SD("access_log", conf.P{"/var/log/nginx/hqs_access_@_www_application_com.log", "json"}, nil),
				conf.SD("error_log", conf.P{"/var/log/nginx/hqs_error_@_www_application_com.log"}, nil),
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
