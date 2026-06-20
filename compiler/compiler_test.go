package compiler

import (
	"testing"

	"gs2parser/parser"
)

func TestCompileBasic(t *testing.T) {
	root, err := parser.Parse(`function onCreated() { temp.a = 1 + 2; }`)
	if err != nil {
		t.Fatal(err)
	}
	code, err := Compile(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) == 0 {
		t.Fatal("empty bytecode")
	}
}
