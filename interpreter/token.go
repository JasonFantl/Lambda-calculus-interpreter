package interpreter

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string // for storing vars
	Position int    // char index
}

const (
	LAMBDA = "\\"
	PERIOD = "."

	FNAME  = "FNAME"
	EQUALS = "="

	LPAREN = "("
	RPAREN = ")"

	VAR = "VAR"

	ILLEGAL = "ILLEGAL"
	NEWLINE = "NEWLINE"
	EOF     = "EOF"
)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

func beginVar(r rune) bool {
	return r >= 'a' && r <= 'z'
}

// allow lowercase, numbers, underscore, and prime
func inVar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '\''
}

func beginFName(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// allow letters, numbers, underscore, and prime
func inFName(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '\''
}
