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
	if p.position >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1] // return last token EOF
	}
	return p.tokens[p.position]
}

func (p *Parser) ParseProgram() (ProgramNode, error) {
	program := ProgramNode{make([]Node, 0)}

	for token := p.currentToken(); token.Type != EOF; token = p.nextToken() {

		expression, err := p.ParseExpression()
		if err != nil {
			return program, err
		}
		program.nodes = append(program.nodes, expression)

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
		varNode, err := p.ParseVar()
		if err != nil {
			return varNode, err
		}
		return varNode, nil
	}

	if token.Type == LPAREN {
		applicationNode, err := p.ParseApplication()
		if err != nil {
			return applicationNode, err
		}
		return applicationNode, nil
	}

	if token.Type == LAMBDA {
		funcNode, err := p.ParseFunction()
		if err != nil {
			return funcNode, err
		}
		return funcNode, nil
	}

	return ProgramNode{}, fmt.Errorf("expression parser expected LAMBDA, APPLICATION, or VAR, got %s at %d", token.Type, token.Position)
}

func (p *Parser) ParseVar() (VarNode, error) {
	varNode := VarNode{}

	token := p.currentToken()

	if token.Type != VAR {
		return varNode, fmt.Errorf("var parser excpected VAR, got %s at %d", token.Type, token.Position)
	}

	varNode.identifier = token.Literal

	return varNode, nil
}

func (p *Parser) ParseApplication() (ApplicationNode, error) {
	applicationNode := ApplicationNode{}

	token := p.currentToken()
	if token.Type != LPAREN {
		return applicationNode, fmt.Errorf("application parser excpected LPAREN, got %s at %d", token.Type, token.Position)
	}

	token = p.nextToken()
	expression, err := p.ParseExpression()
	if err != nil {
		return applicationNode, err
	}
	applicationNode.lExp = expression

	token = p.nextToken()
	expression, err = p.ParseExpression()
	if err != nil {
		return applicationNode, err
	}
	applicationNode.rExp = expression

	token = p.nextToken()
	if token.Type != RPAREN {
		return applicationNode, fmt.Errorf("application parser excpected RPAREN, got %s at %d", token.Type, token.Position)
	}

	return applicationNode, nil
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
		varNode, err := p.ParseVar()
		if err != nil {
			return functionNode, err
		}
		functionNode.input = varNode
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
