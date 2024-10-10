package celpp

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	celast "github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/parser"
)

type Preprocessor struct {
	parser *parser.Parser
}

func New(macros ...cel.Macro) (*Preprocessor, error) {
	pp, err := parser.NewParser(
		parser.Macros(macros...),
	)
	if err != nil {
		return nil, err
	}
	return &Preprocessor{
		parser: pp,
	}, nil
}

func (p *Preprocessor) ProcessToAST(expr string) (*celast.AST, error) {
	ast, iss := p.parser.Parse(common.NewTextSource(expr))
	if len(iss.GetErrors()) > 0 {
		return nil, fmt.Errorf("fail to parse expression: %v", iss.ToDisplayString())
	}
	return ast, nil
}

func (p *Preprocessor) Process(expr string) (string, error) {
	ast, err := p.ProcessToAST(expr)
	if err != nil {
		return "", err
	}
	out, err := parser.Unparse(ast.Expr(), ast.SourceInfo())
	if err != nil {
		return "", fmt.Errorf("unparse: %v", err)
	}
	return out, nil
}
