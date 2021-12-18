package interpreter

type TokenType string

type Token struct {
	Type     TokenType
	Literal  rune // for storing vars
	Position int  // char index
}

const (
	LAMBDA = "\\"
	PERIOD = "."
	LPAREN = "("
	RPAREN = ")"

	VAR = "VAR" // for letters

	ILLEGAL = "ILLEGAL"
	NEWLINE = "NEWLINE"
	EOF     = "EOF"
)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

// later will want to read chars to check for words (pg 16-17).
func isVar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
