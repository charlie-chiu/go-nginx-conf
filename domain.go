package go_nginx_conf

// Config represent nginx Configuration File
type Config struct {
	Directives Block
}

type Block []Directive

// how about least_conn; -> a directiveInterface without params

type Directive interface {
	GetName() string //the directive name.
	GetParameters() []string
	GetComment() []string
	GetBlock() Block
}

type SimpleDirective struct {
	Name            string
	Params, Comment []string
}

func (d SimpleDirective) GetBlock() Block {
	return nil
}

func (d SimpleDirective) GetName() string {
	return d.Name
}

func (d SimpleDirective) GetParameters() []string {
	return d.Params
}

func (d SimpleDirective) GetComment() []string {
	return d.Comment
}

type BlockDirective struct {
	Name            string
	Params, Comment []string
	Block           Block
}

func (d BlockDirective) GetName() string {
	return d.Name
}

func (d BlockDirective) GetParameters() []string {
	return d.Params
}

func (d BlockDirective) GetComment() []string {
	return d.Comment
}

func (d BlockDirective) GetBlock() Block {
	return d.Block
}
