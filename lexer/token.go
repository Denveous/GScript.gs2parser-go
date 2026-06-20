package lexer

type Kind int

const (
	EOF Kind = iota
	Illegal
	Ident
	Int
	Float
	String
	Punct
	Keyword
)

type Token struct {
	Kind Kind
	Lit  string
	Line int
	Col  int
}

func (t Token) Is(s string) bool { return t.Lit == s }
