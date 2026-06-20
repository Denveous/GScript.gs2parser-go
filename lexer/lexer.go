package lexer

import (
	"strconv"
	"strings"
	"unicode"
)

var keywords = map[string]bool{
	"public": true, "if": true, "else": true, "elseif": true, "for": true, "while": true,
	"break": true, "continue": true, "return": true, "function": true, "new": true, "with": true,
	"switch": true, "case": true, "default": true, "const": true, "enum": true, "int": true,
	"float": true, "in": true, "_": true,
}

type Lexer struct {
	src       []rune
	pos       int
	line, col int
}

func Lex(src string) ([]Token, error) {
	l := &Lexer{src: []rune(src), line: 1}
	var out []Token
	for {
		t := l.next()
		out = append(out, t)
		if t.Kind == EOF || t.Kind == Illegal {
			break
		}
	}
	if out[len(out)-1].Kind == Illegal {
		return out, strconv.ErrSyntax
	}
	return out, nil
}

func (l *Lexer) next() Token {
	l.skip()
	startLine, startCol := l.line, l.col
	if l.pos >= len(l.src) {
		return Token{Kind: EOF, Line: startLine, Col: startCol}
	}
	ch := l.peek()
	if isAlpha(ch) {
		s := l.takeWhile(func(r rune) bool { return isAlphaNum(r) || r == ':' })
		if s == "NL" {
			return Token{Kind: Punct, Lit: "@\n", Line: startLine, Col: startCol}
		}
		if s == "SPC" {
			return Token{Kind: Punct, Lit: "@ ", Line: startLine, Col: startCol}
		}
		if s == "TAB" {
			return Token{Kind: Punct, Lit: "@\t", Line: startLine, Col: startCol}
		}
		if keywords[s] {
			return Token{Kind: Keyword, Lit: s, Line: startLine, Col: startCol}
		}
		return Token{Kind: Ident, Lit: s, Line: startLine, Col: startCol}
	}
	if unicode.IsDigit(ch) || (ch == '.' && l.pos+1 < len(l.src) && unicode.IsDigit(l.src[l.pos+1])) {
		return l.number(startLine, startCol)
	}
	if ch == '"' || ch == '\'' {
		return l.string(ch, startLine, startCol)
	}
	for _, op := range []string{"<<=", ">>=", "&&", "||", "==", "!=", "<>", "<=", "=<", ">=", "=>", ":=", "+=", "-=", "*=", "/=", "^=", "%=", "@=", "|=", "&=", "--", "++", "<<", ">>"} {
		if l.has(op) {
			l.advanceN(len([]rune(op)))
			if op == ":=" {
				op = "="
			}
			if op == "=<" {
				op = "<="
			}
			if op == "=>" {
				op = ">="
			}
			if op == "<>" {
				op = "!="
			}
			return Token{Kind: Punct, Lit: op, Line: startLine, Col: startCol}
		}
	}
	l.advance()
	if strings.ContainsRune(".,:;|&(){}[]?!<>+-*/^%=@~", ch) {
		return Token{Kind: Punct, Lit: string(ch), Line: startLine, Col: startCol}
	}
	return Token{Kind: Illegal, Lit: string(ch), Line: startLine, Col: startCol}
}

func (l *Lexer) skip() {
	for l.pos < len(l.src) {
		if l.has("//") {
			for l.pos < len(l.src) && l.peek() != '\n' {
				l.advance()
			}
			continue
		}
		if l.has("/*") {
			l.advanceN(2)
			for l.pos < len(l.src) && !l.has("*/") {
				l.advance()
			}
			if l.has("*/") {
				l.advanceN(2)
			}
			continue
		}
		if l.peek() == ' ' || l.peek() == '\t' || l.peek() == '\r' || l.peek() == '\n' {
			l.advance()
			continue
		}
		break
	}
}

func (l *Lexer) number(line, col int) Token {
	if l.has("0x") || l.has("0X") {
		l.advanceN(2)
		s := l.takeWhile(func(r rune) bool { return unicode.IsDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') })
		v, _ := strconv.ParseInt(s, 16, 32)
		return Token{Kind: Int, Lit: strconv.Itoa(int(v)), Line: line, Col: col}
	}
	s := l.takeWhile(func(r rune) bool { return unicode.IsDigit(r) })
	if l.pos < len(l.src) && l.peek() == '.' {
		s += string(l.advance())
		s += l.takeWhile(func(r rune) bool { return unicode.IsDigit(r) })
		return Token{Kind: Float, Lit: s, Line: line, Col: col}
	}
	return Token{Kind: Int, Lit: s, Line: line, Col: col}
}

func (l *Lexer) string(quote rune, line, col int) Token {
	l.advance()
	var b strings.Builder
	for l.pos < len(l.src) && l.peek() != quote {
		r := l.advance()
		if r == '\\' && l.pos < len(l.src) {
			n := l.advance()
			switch n {
			case 'n':
				b.WriteRune('\n')
			case 't':
				b.WriteRune('\t')
			case 'r':
				b.WriteRune('\r')
			default:
				b.WriteRune(n)
			}
			continue
		}
		b.WriteRune(r)
	}
	if l.pos < len(l.src) {
		l.advance()
	}
	return Token{Kind: String, Lit: b.String(), Line: line, Col: col}
}

func (l *Lexer) takeWhile(fn func(rune) bool) string {
	var b strings.Builder
	for l.pos < len(l.src) && fn(l.peek()) {
		b.WriteRune(l.advance())
	}
	return b.String()
}
func (l *Lexer) has(s string) bool {
	rs := []rune(s)
	if l.pos+len(rs) > len(l.src) {
		return false
	}
	for i, r := range rs {
		if l.src[l.pos+i] != r {
			return false
		}
	}
	return true
}
func (l *Lexer) peek() rune { return l.src[l.pos] }
func (l *Lexer) advance() rune {
	r := l.src[l.pos]
	l.pos++
	if r == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
	return r
}
func (l *Lexer) advanceN(n int) {
	for i := 0; i < n; i++ {
		l.advance()
	}
}
func isAlpha(r rune) bool    { return unicode.IsLetter(r) || r == '_' || r == '$' }
func isAlphaNum(r rune) bool { return isAlpha(r) || unicode.IsDigit(r) }
