package shortcut

import . "go-nginx-conf"

var (
	Listen443SSLHTTP2 = SimpleDirective{Name: "listen", Params: []string{"443", "ssl", "http2"}}
)

func Upstream(upstream string, servers ...SimpleDirective) BlockDirective {
	ups := make([]DirectiveInterface, len(servers))
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
		Block:   &Block{Directives: ups},
	}
}
