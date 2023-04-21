package go_nginx_conf_test

import (
	"bytes"
	"fmt"
	c "github.com/charlie-chiu/go-nginx-conf"
	"os"
	"testing"
)

func TestGenerateL1Conf(t *testing.T) {
	testFixture := "test-fixture/application.l1.conf"
	expected, err := os.ReadFile(testFixture)
	if err != nil {
		t.Fatalf("failed to read test fixture %q, %v", testFixture, err)
	}

	config := c.Config{
		Directives: c.Block{
			c.Upstream("lea_@_www_application_com_443",
				c.UpstreamServer(c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}),
				c.UpstreamServer(c.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"}),
			),
			c.Server(
				c.Listen443SSLHTTP2,
				c.D("server_name", c.P{"www.application.com"}, nil),
				c.D("ssl_stapling", c.P{"on"}, nil),
				c.D("ssl_stapling_verify", c.P{"on"}, nil),
				c.D("ssl_certificate", c.P{"/path/to/cert"}, nil),
				c.D("ssl_certificate_key", c.P{"/path/to/key"}, nil),
				c.D("ssl_protocols", c.P{"TLSv1", "TLSv1.1", "TLSv1.2"}, nil),

				c.D("underscores_in_headers", c.P{"on"}, nil),

				c.ProxySetHeader(c.P{"X-Real-IP", "$remote_addr"}),
				c.ProxySetHeader(c.P{"X-Forwarded-For", "$remote_addr"}),
				c.ProxySetHeader(c.P{"X-Client-Verify", "SUCCESS"}),
				c.ProxySetHeader(c.P{"X-SSL-Subject", "$ssl_client_s_dn"}),
				c.ProxySetHeader(c.P{"X-SSL-Issuer", "$ssl_client_i_dn"}),

				c.D("proxy_http_version", c.P{"1.1"}, nil),
				c.D("proxy_connect_timeout", c.P{"10"}, nil),

				c.If("$http_user_agent ~* \"JianKongBao Monitor\"", c.Return(c.P{"200"})),

				c.D("error_page", c.P{"497", "=307", "https://$host:$server_port$request_uri"}, nil),
				c.D("error_page", c.P{"400", "414", "406", "@requestErrEvent"}, nil),

				c.Location(c.P{"@requestErrEvent"},
					c.D("root", c.P{"/var/tmp/leadns_errpage"}, nil),
					c.D("ssi", c.P{"on"}, nil),
					c.D("internal", nil, nil),
					c.D("try_files", c.P{"$uri", "/custom_error.html", "400"}, nil),
				),
				c.Location(c.P{"~", "/purge(/.*)"},
					c.D("proxy_cache_purge", c.P{"hqszone", "$host$1$is_args$args"}, nil),
				),

				c.D("set", c.P{"$origin_str", "34.96.119.139"}, nil),

				c.Location(c.P{"/"},
					c.D("limit_rate", c.P{"68608k"}, nil),
					c.D("proxy_cache_bypass", c.P{"$nocache"}, nil),
					c.D("proxy_no_cache", c.P{"$nocache"}, nil),
					c.D("proxy_cache", c.P{"hqszone"}, nil),
					c.D("proxy_cache_valid", c.P{"301", "302", "0s"}, nil),
					c.D("proxy_cache_key", c.P{"$host$uri$is_args$args$origin_str"}, nil),
					c.D("proxy_pass", c.P{"$scheme://lea_@_www_application_com_443"}, nil),
				),

				c.D("access_log", c.P{"/var/log/nginx/hqs_access_@_www_application_com.log", "json"}, nil),
				c.D("error_log", c.P{"/var/log/nginx/hqs_error_@_www_application_com.log"}, nil),
			),
		},
	}

	actual := c.DumpConfig(config, c.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateSimpleDirective(t *testing.T) {
	type testCase struct {
		name   string
		input  c.SimpleDirective
		output []byte
	}

	testCases := []testCase{
		{
			name: "listen 80",
			input: c.SimpleDirective{
				Name:   "listen",
				Params: c.P{"80"},
			},
			output: []byte("listen 80;"),
		},
		{
			name: "custom error page",
			input: c.SimpleDirective{
				Name:   "error_page",
				Params: c.P{"497", "=307", "https://$host:$server_port$request_uri"},
			},
			output: []byte("error_page 497 =307 https://$host:$server_port$request_uri;"),
		},
		{
			name: "set variable",
			input: c.SimpleDirective{
				Name:   "set",
				Params: c.P{"$origin_str", "34.96.119.139"},
			},
			output: []byte("set $origin_str 34.96.119.139;"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test simple directive: %q", tc.name), func(t *testing.T) {
			actual := c.DumpConfig(c.Config{
				Directives: c.Block{
					tc.input,
				},
			}, c.IndentedStyle)

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

	config := c.Config{
		Directives: c.Block{
			c.Upstream("lea_@_www_application_com_443",
				c.SimpleDirective{
					Name:   "server",
					Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"},
				},
				c.SimpleDirective{
					Name:   "server",
					Params: c.P{"34.92.95.215:80", "max_fails=1", "fail_timeout=10s"},
				},
			)},
	}
	actual := c.DumpConfig(config, c.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateIfDirective(t *testing.T) {
	expected := []byte(
		`if ($host = 'www.application.com') {
    return 301 https://$host$request_uri;
}`)

	config := c.Config{
		Directives: c.Block{
			c.If(
				"$host = 'www.application.com'",
				c.SimpleDirective{
					Name:   "return",
					Params: c.P{"301", "https://$host$request_uri"},
				},
			),
		},
	}

	actual := c.DumpConfig(config, c.IndentedStyle)

	assertConfigEqual(t, expected, actual)
}

func TestGenerateLocationDirective(t *testing.T) {
	expected := []byte(`location ~ /purge(/.*) {
    proxy_cache_purge hqszone $host$1$is_args$args;
}`)

	config := c.Config{
		Directives: c.Block{
			c.Location(
				c.P{"~ /purge(/.*)"},
				c.SimpleDirective{
					Name: "proxy_cache_purge", Params: c.P{"hqszone", "$host$1$is_args$args"},
				},
			),
		},
	}

	actual := c.DumpConfig(config, c.IndentedStyle)

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
