package interpreter

import (
	"fmt"
)

type Parser struct {
	tokens   []Token
	position int // token index
	root     ProgramNode
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
		root:     ProgramNode{make([]Node, 0)},
	}
}

func (p *Parser) Parse() (ProgramNode, error) {
	return p.ParseProgram()
}

func (p *Parser) nextToken() Token {
	p.position++
	return p.currentToken()
}

func (p *Parser) currentToken() Token {
	if p.position < len(p.tokens) {
		return p.tokens[p.position]
	}
	return p.tokens[len(p.tokens)-1] // return last token EOF
}

func (p *Parser) ParseProgram() (ProgramNode, error) {
	program := ProgramNode{make([]Node, 0)}

	for token := p.currentToken(); token.Type != EOF; token = p.nextToken() {

		// ignore all NEWLINEs before an expression
		if token.Type == NEWLINE {
			continue
		}

		// if FNAME, try def, otherwise go to exp to try call. FIX THIS
		savedTokenIndex := p.position
		wasFuncDef := false
		if token.Type == FNAME {
			namedFunctionDefNode, err := p.ParseNamedFunctionDef()
			if err == nil { // was a func def!
				program.nodes = append(program.nodes, namedFunctionDefNode)
				wasFuncDef = true
			}
		}

		if !wasFuncDef {
			p.position = savedTokenIndex // in case of failed func def parse
			expression, err := p.ParseExpression()
			if err != nil {
				return program, err
			}
			program.nodes = append(program.nodes, expression)
		}

		token = p.nextToken()
		if token.Type != NEWLINE && token.Type != EOF {
			return program, fmt.Errorf("program parser expected NEWLINE or EOF, got %s at %d", token.Type, token.Position)
		}
	}

	return program, nil
}

func (p *Parser) ParseExpression() (Node, error) {
	token := p.currentToken()

	if token.Type == VAR {
		return VarNode{token.Literal}, nil
	}

	if token.Type == LAMBDA {
		funcNode, err := p.ParseFunction()
		if err != nil {
			return funcNode, err
		}
		return funcNode, nil
	}

	if token.Type == LPAREN { // could be three different things
		token = p.nextToken()
		if token.Type == FNAME { // parse function call
			// back up token position since parser requires a starting LPAREN
			p.position--
			fCallNode, err := p.ParseNamedFunctionCall()
			if err != nil {
				return fCallNode, err
			}
			return fCallNode, nil
		} else { // narrowed down to two
			// parse the (first?) expression
			expression1, err := p.ParseExpression()
			if err != nil {
				return ProgramNode{}, err
			}

			token = p.nextToken()
			if token.Type == RPAREN { // determined single var
				return expression1, nil
			} else { // must be application
				expression2, err := p.ParseExpression()
				if err != nil {
					return ProgramNode{}, err
				}

				token = p.nextToken()
				if token.Type != RPAREN {
					return ProgramNode{}, fmt.Errorf("application parser excpected RPAREN, got %s at %d", token.Type, token.Position)
				}

				return ApplicationNode{expression1, expression2}, nil
			}
		}
	}

	return ProgramNode{}, fmt.Errorf("expression parser expected LAMBDA, APPLICATION, FNAME, or VAR, got %s at %d", token.Type, token.Position)
}

func (p *Parser) ParseNamedFunctionDef() (NamedFunctionDefNode, error) {
	namedFunctionDefNode := NamedFunctionDefNode{}

	token := p.currentToken()
	if token.Type != FNAME {
		return namedFunctionDefNode, fmt.Errorf("named function def parser excpected FNAME, got %s at %d", token.Type, token.Position)
	} else {
		namedFunctionDefNode.identifier = token.Literal
	}

	// now capture the rest
	for token = p.nextToken(); token.Type == VAR; token = p.nextToken() {
		namedFunctionDefNode.inputs = append(namedFunctionDefNode.inputs, VarNode{token.Literal})
	}

	// token = p.nextToken() // we overshot in last check, next token already loaded
	if token.Type != EQUALS {
		return namedFunctionDefNode, fmt.Errorf("named function def parser excpected EQUALS, got %s at %d", token.Type, token.Position)
	}

	token = p.nextToken()
	expression, err := p.ParseExpression()
	if err != nil {
		return namedFunctionDefNode, err
	}
	namedFunctionDefNode.body = expression

	return namedFunctionDefNode, nil
}

func (p *Parser) ParseNamedFunctionCall() (NamedFunctionCallNode, error) {
	namedFunctionCallNode := NamedFunctionCallNode{}

	token := p.currentToken()
	if token.Type != LPAREN {
		return namedFunctionCallNode, fmt.Errorf("named function call parser excpected LPAREN, got %s at %d", token.Type, token.Position)
	}

	token = p.nextToken()
	if token.Type != FNAME {
		return namedFunctionCallNode, fmt.Errorf("named function call parser excpected FNAME, got %s at %d", token.Type, token.Position)
	}
	namedFunctionCallNode.identifier = token.Literal

	// capture all the EXP
	for token = p.nextToken(); token.Type != RPAREN; token = p.nextToken() {
		expression, err := p.ParseExpression()
		if err != nil {
			return namedFunctionCallNode, err
		}

		namedFunctionCallNode.inputs = append(namedFunctionCallNode.inputs, expression)
	}

	// token = p.nextToken() // we overshot in last check, next token already loaded
	if token.Type != RPAREN {
		return namedFunctionCallNode, fmt.Errorf("named function call parser excpected RPAREN, got %s at %d", token.Type, token.Position)
	}

	return namedFunctionCallNode, nil
}

func (p *Parser) ParseFunction() (FunctionNode, error) {
	functionNode := FunctionNode{}

	token := p.currentToken()
	if token.Type != LAMBDA {
		return functionNode, fmt.Errorf("function parser excpected LAMBDA, got %s at %d", token.Type, token.Position)
	}

	token = p.nextToken()
	if token.Type != VAR {
		return functionNode, fmt.Errorf("function parser excpected VAR, got %s at %d", token.Type, token.Position)
	} else {
		functionNode.input = VarNode{token.Literal}
	}

	token = p.nextToken()
	if token.Type != PERIOD {
		return functionNode, fmt.Errorf("function parser excpected PERIOD, got %s at %d", token.Type, token.Position)
	}

	token = p.nextToken()
	expression, err := p.ParseExpression()
	if err != nil {
		return functionNode, err
	}
	functionNode.body = expression

	return functionNode, nil
}
