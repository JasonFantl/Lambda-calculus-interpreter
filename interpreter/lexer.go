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
	for err == nil && t.Type != EOF_TOKEN {
		// ignore illegal tokens
		if t.Type == ILLEGAL_TOKEN {
			fmt.Printf("illegal token:\n\t%s\t%*s\n", string(l.input), t.Position+1, "^")
		} else {
			tokens = append(tokens, t)
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
	return l.currentRune()
}

func (l *Lexer) currentRune() (rune, bool) {
	if l.position < len(l.input) {
		return l.input[l.position], false
	}
	return 0, true
}

var tokenLookup = map[string]TokenType{
	"\\":  LAMBDA_TOKEN,
	".":   PERIOD_TOKEN,
	"(":   LPAREN_TOKEN,
	")":   RPAREN_TOKEN,
	"\n":  NEWLINE_TOKEN,
	"def": DEF_TOKEN,
}

func (l *Lexer) nextToken() (Token, error) {

	r, eof := l.nextRune()

	// handle whitespace (maybe use later?)
	for !eof && isWhitespace(r) {
		r, eof = l.nextRune()
	}

	token := Token{EOF_TOKEN, "", l.position}
	if eof {
		return token, nil
	}

	if val, ok := tokenLookup[string(r)]; ok {
		token.Type = val
	} else {
		word := l.lexWord()

		// check keywords
		if val, ok := tokenLookup[word]; ok {
			token.Type = val
		} else if isVar(word) {
			token.Type = VAR_TOKEN
			token.Literal = word
		} else if isName(word) {
			token.Type = NAME_TOKEN
			token.Literal = word
		} else {
			return token, fmt.Errorf("illegal character '%s' at %d", string(r), token.Position)
		}
	}

	return token, nil
}

func (l *Lexer) lexWord() string {
	s := ""
	for r, eof := l.currentRune(); inWord(r) && !eof; {
		s += string(r)
		r, eof = l.nextRune()
	}
	l.position-- // we overshot by a rune to check

	return s
}

// var begin with lowercase letter
func isVar(s string) bool {
	// check first rune
	b := s[0]
	if !(b >= 'a' && b <= 'z') {
		return false
	}
	// check following runes
	for _, r := range s {
		if !inWord(r) {
			return false
		}
	}
	return true
}

// names begin with uppercase letter
func isName(s string) bool {
	// check first rune
	b := s[0]
	if !(b >= 'A' && b <= 'Z') {
		return false
	}
	// check following runes
	for _, r := range s {
		if !inWord(r) {
			return false
		}
	}
	return true
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

func inWord(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_' ||
		r == '\''
}
