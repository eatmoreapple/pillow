package internal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eatmoreapple/juice"
)

type Generator struct {
	cfg       *juice.Configuration
	impl      *Implement
	namespace string
	args      string
	output    string
	writer    io.Writer
}

func (g *Generator) Generate() error {
	for _, method := range g.impl.Methods {
		key := fmt.Sprintf("%s.%s", g.namespace, method.Name)
		statement, err := g.cfg.Mappers.GetStatementByID(key)
		if err != nil {
			return err
		}
		maker := FunctionBodyMaker{statement: statement, function: method}
		if err := maker.Make(); err != nil {
			return err
		}
	}
	g.impl.ExtraImports = append(g.impl.ExtraImports, Import{Path: "github.com/eatmoreapple/juice", Name: "juice"})
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("// Code generated by \"%s\"; DO NOT EDIT.", g.args))
	builder.WriteString("\n\n")
	builder.WriteString(g.impl.String())
	writer, err := g.Output()
	if err != nil {
		return err
	}
	defer func() {
		if closer, ok := writer.(io.Closer); ok {
			_ = closer.Close()
		}
	}()
	_, err = io.WriteString(writer, builder.String())
	return err
}

func (g *Generator) Output() (_ io.Writer, err error) {
	if g.writer == nil {
		if len(g.output) > 0 {
			if !strings.HasSuffix(g.output, ".go") {
				g.output += ".go"
			}
			g.writer, err = os.Create(g.output)
			if err != nil {
				return nil, err
			}
		} else {
			g.writer = os.Stdout
		}
	}
	return g.writer, nil
}
