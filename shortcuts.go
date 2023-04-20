package go_nginx_conf

import (
	"fmt"
)

type P []string

var (
	Listen80          = SimpleDirective{Name: "listen", Params: []string{"80"}}
	Listen443SSLHTTP2 = SimpleDirective{Name: "listen", Params: []string{"443", "ssl", "http2"}}
)

// D is stand for simple directive
func D(name string, params, comment []string) SimpleDirective {
	return SimpleDirective{
		Name:    name,
		Params:  params,
		Comment: comment,
	}
}

func Upstream(upstream string, servers ...SimpleDirective) BlockDirective {
	ups := make(Block, len(servers))
	for i, server := range servers {
		ups[i] = SimpleDirective{
			Name:   server.Name,
			Params: server.Params,
		}
	}
	ups = append(ups, SimpleDirective{
		Name:    "least_conn",
		Params:  nil,
		Comment: nil,
	})

	return BlockDirective{
		Name:    "upstream",
		Params:  []string{upstream},
		Comment: nil,
		Block:   ups,
	}
}

func UpstreamServer(params []string) SimpleDirective {
	return SimpleDirective{
		Name:   "server",
		Params: params,
	}
}

func Server(directives ...Directive) BlockDirective {
	return BlockDirective{
		Name:  "server",
		Block: directives,
	}
}

func Location(parameters []string, directives ...Directive) BlockDirective {
	return BlockDirective{
		Name:   "location",
		Params: parameters,
		Block:  directives,
	}
}

func If(condition string, directives ...Directive) BlockDirective {
	return BlockDirective{
		Name:   "if",
		Params: []string{fmt.Sprintf("(%s)", condition)},
		Block:  directives,
	}
}

func ProxySetHeader(params []string) SimpleDirective {
	return SimpleDirective{
		Name:   "proxy_set_header",
		Params: params,
	}
}

func Return(params []string) SimpleDirective {
	return SimpleDirective{
		Name:   "return",
		Params: params,
	}
}
