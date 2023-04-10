package go_nginx_conf_test

import (
	"bytes"
	go_nginx_conf "go-nginx-conf"
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

	config := go_nginx_conf.Config{
		Directives: &go_nginx_conf.Block{Directives: []go_nginx_conf.DirectiveInterface{
			c.Upstream("lea_@_www_jb1228_com_80",
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

func assertConfigEqual(t *testing.T, expected []byte, actual []byte) {
	t.Helper()
	if bytes.Compare(expected, actual) != 0 {
		t.Logf("\nfailed to assert actual equal expected\n")
		t.Logf("expected:\n%s\n", expected)
		t.Logf("actual:  \n%s\n", actual)
		t.Fail()
	}
}
