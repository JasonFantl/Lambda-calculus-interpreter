package interpreter

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string // for storing vars
	Position int    // char index
}

func (t Token) String() string {
	return string(t.Type)
}

const (
	LAMBDA_TOKEN = "\\"
	PERIOD_TOKEN = "."

	NAME_TOKEN = "NAME"
	DEF_TOKEN  = "DEF"

	LPAREN_TOKEN = "("
	RPAREN_TOKEN = ")"

	VAR_TOKEN = "VAR"

	ILLEGAL_TOKEN = "ILLEGAL"
	NEWLINE_TOKEN = "NEWLINE"
	EOF_TOKEN     = "EOF"
)
