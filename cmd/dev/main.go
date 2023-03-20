package main

import (
	"fmt"
	nc "go-nginx-conf"
)

func main() {
	directives := []nc.DirectiveInterface{
		nc.SimpleDirective{
			Name:    "least_conn",
			Params:  nil,
			Comment: nil,
		},
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
