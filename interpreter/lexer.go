package interpreter

import (
	"fmt"
)

type Lexer struct {
	input    []rune
	position int // ccurrent runes index
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    []rune(input),
		position: -1, // when we call next rune all is good
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0)

	t, err := l.nextToken()
	for err == nil && t.Type != EOF {
		// ignore illegal tokens
		if t.Type != ILLEGAL {
			tokens = append(tokens, t)
		} else {
			fmt.Printf("illegal token:\n\t%s\t%*s\n", string(l.input), t.Position+1, "^")
		}
		t, err = l.nextToken()
	}

	if err != nil {
		return tokens, err
	}

	// if we exited not due to error, must have been EOF
	tokens = append(tokens, t)

	return tokens, err
}

func (l *Lexer) nextRune() (rune, bool) {
	l.position++
	if l.position < len(l.input) {
		return l.input[l.position], false
	}
	return 0, true
}

func (l *Lexer) nextToken() (Token, error) {

	// since every token is a rune, no need to keep track of words (yet)
	r, eof := l.nextRune()

	// handle whitespace (maybe use later?)
	for !eof && isWhitespace(r) {
		r, eof = l.nextRune()
	}

	token := Token{EOF, 0, l.position}
	if eof {
		return token, nil
	}

	switch r {
	case '\\':
		token.Type = LAMBDA
	case '.':
		token.Type = PERIOD
	case '(':
		token.Type = LPAREN
	case ')':
		token.Type = RPAREN
	case '\n':
		token.Type = NEWLINE
	default:
		if isVar(r) {
			token.Type = VAR
			token.Literal = r
		} else {
			token.Type = ILLEGAL
		}
	}

	return token, nil
}
