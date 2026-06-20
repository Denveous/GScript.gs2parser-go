package main

import (
	"flag"
	"fmt"
	"os"

	"gs2parser"
)

func main() {
	out := flag.String("o", "", "output bytecode file")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: gs2compiler [-o out.bytecode] script.gs2")
		os.Exit(2)
	}
	src, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		die(err)
	}
	res, err := gs2parser.Compile(string(src))
	if err != nil {
		die(err)
	}
	if *out == "" {
		os.Stdout.Write(res.Bytecode)
		return
	}
	if err := os.WriteFile(*out, res.Bytecode, 0644); err != nil {
		die(err)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
