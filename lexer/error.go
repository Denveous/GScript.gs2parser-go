package lexer

import "fmt"

type Error struct {
	Message string
	Line    int
	Column  int
	Near    string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	msg := e.Message
	if e.Line > 0 {
		msg += fmt.Sprintf(" at %d:%d", e.Line, e.Column)
	}
	if e.Near != "" {
		msg += fmt.Sprintf(" near %q", e.Near)
	}
	return msg
}
