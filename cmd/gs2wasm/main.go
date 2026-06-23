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
			return result(nil, "missing source argument", nil)
		}
		res := gs2parser.CompileDetailed(args[0].String())
		if len(res.Diagnostics) != 0 {
			return result(nil, (&gs2parser.DiagnosticError{Diagnostics: res.Diagnostics}).Error(), res.Diagnostics)
		}
		return result(res.Bytecode, "", nil)
	}))
	<-done
}

func result(code []byte, msg string, diagnostics []gs2parser.Diagnostic) js.Value {
	obj := js.Global().Get("Object").New()
	if msg != "" {
		obj.Set("ok", false)
		obj.Set("error", msg)
		if len(diagnostics) != 0 {
			obj.Set("diagnostics", diagnosticsValue(diagnostics))
		}
		return obj
	}
	arr := js.Global().Get("Uint8Array").New(len(code))
	js.CopyBytesToJS(arr, code)
	obj.Set("ok", true)
	obj.Set("bytecode", arr)
	return obj
}

func diagnosticsValue(diagnostics []gs2parser.Diagnostic) js.Value {
	arr := js.Global().Get("Array").New()
	for _, d := range diagnostics {
		obj := js.Global().Get("Object").New()
		obj.Set("severity", d.Severity)
		obj.Set("stage", d.Stage)
		obj.Set("message", d.Message)
		obj.Set("line", d.Line)
		obj.Set("column", d.Column)
		obj.Set("near", d.Near)
		obj.Set("sourceLine", d.SourceLine)
		arr.Call("push", obj)
	}
	return arr
}
