package gs2parser

import (
	"errors"
	"fmt"
	"strings"

	"gs2parser/lexer"
	"gs2parser/parser"
)

type Diagnostic struct {
	Severity   string
	Stage      string
	Message    string
	Line       int
	Column     int
	Near       string
	SourceLine string
}

type DiagnosticError struct {
	Diagnostics []Diagnostic
}

func (e *DiagnosticError) Error() string {
	if e == nil || len(e.Diagnostics) == 0 {
		return ""
	}
	if len(e.Diagnostics) == 1 {
		return e.Diagnostics[0].Error()
	}
	return fmt.Sprintf("%s (+%d more)", e.Diagnostics[0].Error(), len(e.Diagnostics)-1)
}

func (d Diagnostic) Error() string {
	msg := d.Message
	if d.Line > 0 {
		msg += fmt.Sprintf(" at %d:%d", d.Line, d.Column)
	}
	if d.Near != "" {
		msg += fmt.Sprintf(" near %q", d.Near)
	}
	return msg
}

func diagnosticsFromError(source, stage string, err error) []Diagnostic {
	var lexErr *lexer.Error
	if errors.As(err, &lexErr) {
		return []Diagnostic{{Severity: "error", Stage: "lexer", Message: lexErr.Message, Line: lexErr.Line, Column: lexErr.Column, Near: lexErr.Near, SourceLine: sourceLine(source, lexErr.Line)}}
	}
	var parseErr *parser.Error
	if errors.As(err, &parseErr) {
		return []Diagnostic{{Severity: "error", Stage: "parser", Message: parseErr.Message, Line: parseErr.Line, Column: parseErr.Column, Near: parseErr.Near, SourceLine: sourceLine(source, parseErr.Line)}}
	}
	return []Diagnostic{{Severity: "error", Stage: stage, Message: err.Error()}}
}

func sourceLine(source string, line int) string {
	if line < 1 {
		return ""
	}
	lines := strings.Split(source, "\n")
	if line > len(lines) {
		return ""
	}
	return strings.TrimSuffix(lines[line-1], "\r")
}
