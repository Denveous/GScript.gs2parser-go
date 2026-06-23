package gs2parser

import (
	"github.com/MorenoLand/GScript.gs2parser-go/ast"
	"github.com/MorenoLand/GScript.gs2parser-go/compiler"
	"github.com/MorenoLand/GScript.gs2parser-go/parser"
)

type Result struct {
	Bytecode    []byte
	AST         *ast.Block
	Diagnostics []Diagnostic
}

func Parse(source string) (*ast.Block, error) {
	return parser.Parse(source)
}

func Compile(source string) (*Result, error) {
	res := CompileDetailed(source)
	if len(res.Diagnostics) != 0 {
		return nil, &DiagnosticError{Diagnostics: res.Diagnostics}
	}
	return res, nil
}

func CompileDetailed(source string) *Result {
	res := &Result{}
	root, err := parser.Parse(source)
	if err != nil {
		res.Diagnostics = diagnosticsFromError(source, "parser", err)
		return res
	}
	res.AST = root
	code, err := compiler.Compile(root)
	if err != nil {
		res.Diagnostics = diagnosticsFromError(source, "compiler", err)
		return res
	}
	res.Bytecode = code
	return res
}
