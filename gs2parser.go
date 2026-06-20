package gs2parser

import (
	"gs2parser/ast"
	"gs2parser/compiler"
	"gs2parser/parser"
)

type Result struct {
	Bytecode []byte
	AST      *ast.Block
}

func Parse(source string) (*ast.Block, error) {
	return parser.Parse(source)
}

func Compile(source string) (*Result, error) {
	root, err := parser.Parse(source)
	if err != nil {
		return nil, err
	}
	code, err := compiler.Compile(root)
	if err != nil {
		return nil, err
	}
	return &Result{Bytecode: code, AST: root}, nil
}
