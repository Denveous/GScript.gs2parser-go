//go:build js && wasm

package main

import (
	"syscall/js"

	"gs2parser"
)

func main() {
	done := make(chan struct{})
	js.Global().Set("gs2Compile", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return result(nil, "missing source argument")
		}
		res, err := gs2parser.Compile(args[0].String())
		if err != nil {
			return result(nil, err.Error())
		}
		return result(res.Bytecode, "")
	}))
	<-done
}

func result(code []byte, msg string) js.Value {
	obj := js.Global().Get("Object").New()
	if msg != "" {
		obj.Set("ok", false)
		obj.Set("error", msg)
		return obj
	}
	arr := js.Global().Get("Uint8Array").New(len(code))
	js.CopyBytesToJS(arr, code)
	obj.Set("ok", true)
	obj.Set("bytecode", arr)
	return obj
}
