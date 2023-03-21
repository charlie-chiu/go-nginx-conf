package main

import (
	"fmt"
	nc "go-nginx-conf"
	c "go-nginx-conf/shortcut"
)

func main() {
	config := nc.Config{
		Directives: &nc.Block{Directives: []nc.DirectiveInterface{
			c.Upstream("lea_@_www_jb1228_com_80",
				nc.SimpleDirective{Name: "server", Params: c.P{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}},
			),
			c.Upstream("lea_@_www_jb1228_com_443",
				nc.SimpleDirective{Name: "server", Params: c.P{"35.200.43.88:443", "max_fails=1", "fail_timeout=10s"}},
			),
			c.Server(
				c.Listen80,
				c.Location(c.P{"@relayEvent"}, nc.SimpleDirective{Name: "try_files", Params: c.P{"$uri", "/custom_error.html", "403"}}),
				c.Location(c.P{"/"}, nc.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_jb1228_com_443"}}),
			),
			c.Server(
				c.Listen443SSLHTTP2,
				c.Location(c.P{"/"}, nc.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_jb1228_com_443"}}),
				c.Location(c.P{"/"}, nc.SimpleDirective{Name: "proxy_pass", Params: c.P{"$scheme://lea_@_www_jb1228_com_443"}}),
			),
		}},
	}

	fmt.Printf("%s\n", nc.DumpConfig(config, nc.IndentedStyle))
}
