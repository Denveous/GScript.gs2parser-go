package parser

import "testing"

func TestParseBasic(t *testing.T) {
	_, err := Parse(`function onCreated() { temp.a = 1 + 2; echo("ok"); }`)
	if err != nil {
		t.Fatal(err)
	}
}
