package main

import (
	"fmt"
	nc "go-nginx-conf"
	c "go-nginx-conf/shortcut"
)

func main() {
	directives := []nc.DirectiveInterface{
		c.Listen443SSLHTTP2,
		c.Upstream(
			"lea_@_www_jb1228_com_80",
			nc.SimpleDirective{Name: "server", Params: []string{"35.200.43.88:80", "max_fails=1", "fail_timeout=10s"}},
		),
		nc.SimpleDirective{
			Name:    "proxy_redirect",
			Params:  []string{"off"},
			Comment: nil,
		},
		nc.SimpleDirective{
			Name:    "proxy_buffers",
			Params:  []string{"4", "32k"},
			Comment: []string{"#with some comment"},
		},
		nc.BlockDirective{
			Name:    "location",
			Params:  []string{"/404.html"},
			Comment: nil,
			Block: &nc.Block{
				Directives: []nc.DirectiveInterface{
					nc.SimpleDirective{
						Name:   "return",
						Params: []string{"200"},
					},
				},
			},
		},
		nc.BlockDirective{
			Name:    "location",
			Params:  nil,
			Comment: nil,
			Block: &nc.Block{
				Directives: []nc.DirectiveInterface{
					nc.SimpleDirective{
						Name:   "return",
						Params: []string{"200"},
					},
				},
			},
		},
	}

	config := nc.Config{
		Directives: &nc.Block{Directives: directives},
	}

	fmt.Printf("%s\n", nc.DumpConfig(config, nc.IndentedStyle))
}
