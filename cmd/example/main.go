package main

import (
	"fmt"
	c "github.com/charlie-chiu/go-nginx-conf"
)

func main() {
	config := c.Config{
		Directives: []c.Directive{
			c.Upstream("lea_@_www_jb1228_com_80",
				c.SimpleDirective{Name: "server", Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}},
			),
			c.Upstream("lea_@_www_jb1228_com_443",
				c.SimpleDirective{Name: "server", Params: c.P{"35.200.43.88:443", "max_fails=1", "fail_timeout=10s"}},
			),
			c.Server(
				c.Listen80,
				c.Location(c.P{"@relayEvent"}, c.SimpleDirective{Name: "try_files", Params: c.P{"$uri", "/custom_error.html", "403"}}),
				c.Location(c.P{"/"}, c.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_jb1228_com_443"}}),
			),
			c.Server(
				c.Listen443SSLHTTP2,
				c.Location(c.P{"/"}, c.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_jb1228_com_443"}}),
				c.Location(c.P{"/"}, c.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_jb1228_com_443"}}),
			),
		},
	}

	fmt.Printf("%s\n", c.DumpConfig(config, c.IndentedStyle))
}
