package go_nginx_conf

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var (
	//NoIndentStyle default style
	NoIndentStyle = &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         0,
	}

	//IndentedStyle default style
	IndentedStyle = &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
	}
)

// Style dumping style
type Style struct {
	SortDirectives    bool
	SpaceBeforeBlocks bool
	StartIndent       int
	Indent            int
}

// NewStyle create new style
func NewStyle() *Style {
	style := &Style{
		SortDirectives: false,
		StartIndent:    0,
		Indent:         4,
	}
	return style
}

// Iterate interate the indentation for sub blocks
func (s *Style) Iterate() *Style {
	newStyle := &Style{
		SortDirectives:    s.SortDirectives,
		SpaceBeforeBlocks: s.SpaceBeforeBlocks,
		StartIndent:       s.StartIndent + s.Indent,
		Indent:            s.Indent,
	}
	return newStyle
}

func DumpConfig(config Config, style *Style) []byte {
	return DumpBlock(config.Directives, style)
}

func DumpDirective(d Directive, style *Style) []byte {
	var buf bytes.Buffer

	if style.SpaceBeforeBlocks && d.GetBlock() != nil {
		buf.WriteString("\n")
	}
	if len(d.GetComment()) > 0 {
		for _, comment := range d.GetComment() {
			buf.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", style.StartIndent), comment))
		}
	}
	buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(" ", style.StartIndent), d.GetName()))
	if len(d.GetParameters()) > 0 {
		buf.WriteString(fmt.Sprintf(" %s", strings.Join(d.GetParameters(), " ")))
	}
	if d.GetBlock() == nil {
		buf.WriteRune(';')
	} else {
		buf.WriteString(" {\n")
		buf.Write(DumpBlock(d.GetBlock(), style.Iterate()))
		buf.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(" ", style.StartIndent)))
	}

	return buf.Bytes()
	//return []byte("dummy output of DumpDirective")
}

func DumpBlock(b Block, style *Style) []byte {
	var buf bytes.Buffer

	directives := b
	if style.SortDirectives {
		sort.SliceStable(directives, func(i, j int) bool {
			return directives[i].GetName() < directives[j].GetName()
		})
	}

	for i, directive := range directives {
		buf.Write(DumpDirective(directive, style))
		if i != len(directives)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.Bytes()
}
