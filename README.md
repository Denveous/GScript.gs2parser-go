# gs2parser

This is a Go compiler for the Graal Script 2 (GS2) language.

It is a Go port of the C++ `gs2-parser` project, structured so the compiler can be used as its own command-line tool or imported as a module by another Go project.

# Prerequisites

Install Go 1.21 or newer.

# Building

You can build the CLI with:

```sh
go build ./cmd/gs2compiler
```

# Building (Wasm)

The compiler core is pure Go and can be built for WebAssembly using Go's `js/wasm` target:

```sh
GOOS=js GOARCH=wasm go build -o gs2parser.wasm ./cmd/gs2wasm
```

The resulting `gs2parser.wasm` file can be imported into a webpage using Go's `wasm_exec.js`. It exposes:

```js
const result = gs2Compile(source);
```

`result.ok` is `true` on success and `result.bytecode` is a `Uint8Array`.

# Running

The non-wasm build can be run using:

```sh
gs2compiler script.gs2
gs2compiler -o script.gs2bc script.gs2
gs2compiler script.gs2 script.gs2bc
gs2compiler scripts/
gs2compiler file1.gs2 file2.gs2 file3.gs2
```

# Library Usage

```go
package main

import "gs2parser"

func main() {
	result, err := gs2parser.Compile(`function onCreated() { temp.a = 1; }`)
	if err != nil {
		panic(err)
	}
	_ = result.Bytecode
}
```
