package interpreter

import (
	"fmt"
)

type EngParser struct {
	tokens   []Token
	position int // token index
	root     ProgramNode
}

func NewEngParser(tokens []Token) *EngParser {
	return &EngParser{
		tokens:   tokens,
		position: 0,
		root:     ProgramNode{make([]Node, 0)},
	}
}

func (p *EngParser) Parse() (ProgramNode, error) {
	return p.ParseProgram()
}

func (p *EngParser) advanceToken() Token {
	p.position++
	return p.currentToken()
}

func (p *EngParser) currentToken() Token {
	if p.position < len(p.tokens) {
		return p.tokens[p.position]
	}
	return p.tokens[len(p.tokens)-1] // return last token EOF
}

// NOTE: each parse ends at the last successful token of the parse

func (p *EngParser) ParseProgram() (ProgramNode, error) {
	program := ProgramNode{make([]Node, 0)}

	for token := p.currentToken(); token.Type != EOF_TOKEN; token = p.advanceToken() {

		// ignore all NEWLINEs before an expression
		if token.Type == NEWLINE_TOKEN {
			continue
		}

		expression, err := p.ParseExp(0)
		if err != nil {
			return program, fmt.Errorf("program parser -> %s", err)
		}
		program.nodes = append(program.nodes, expression)
	}

	return program, nil
}

// https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm#climbing
func (p *EngParser) ParseExp(prec int) (Node, error) {
	exp, err := p.ParseP()
	if err != nil {
		return nil, fmt.Errorf("expression parser -> %s", err)
	}

	// we use FIRST to check if we have an "application operator"
	for nextToken := p.advanceToken(); parseExpressionPrimeFIRST[nextToken.Type] &&
		prec <= 0; nextToken = p.advanceToken() {

		q := prec + 1

		appNode, err := p.ParseExp(q)
		if err != nil {
			return ProgramNode{}, fmt.Errorf("expression parser -> %s", err)
		}

		exp = ApplicationNode{exp, appNode}
	}
	// we went too far, need to decrement
	p.position--

	return exp, nil
}

func (p *EngParser) ParseP() (Node, error) {
	token := p.currentToken()

	if token.Type == VAR_TOKEN {
		varNode, err := p.ParseVar()
		if err != nil {
			return varNode, fmt.Errorf("p parser -> %s", err)
		}
		return varNode, nil
	}

	if token.Type == LAMBDA_TOKEN {
		funcNode, err := p.ParseFunction()
		if err != nil {
			return funcNode, fmt.Errorf("p parser -> %s", err)
		}
		return funcNode, nil
	}

	if token.Type == LPAREN_TOKEN {
		groupedNode, err := p.ParseGrouped()
		if err != nil {
			return groupedNode, fmt.Errorf("p parser -> %s", err)
		}
		return groupedNode, nil
	}

	return ProgramNode{}, fmt.Errorf("p parser -> expected VAR, LAMBDA, or LPAREN, got %s at %d", token.Type, token.Position)
}

func (p *EngParser) ParseVar() (VarNode, error) {
	varNode := VarNode{}

	token := p.currentToken()
	if token.Type != VAR_TOKEN {
		return varNode, fmt.Errorf("var parser -> excpected VAR, got %s at %d", token.Type, token.Position)
	}
	varNode.identifier = token.Literal

	return varNode, nil
}

func (p *EngParser) ParseFunction() (FunctionNode, error) {
	functionNode := FunctionNode{}

	token := p.currentToken()
	if token.Type != LAMBDA_TOKEN {
		return functionNode, fmt.Errorf("function parser -> excpected LAMBDA, got %s at %d", token.Type, token.Position)
	}

	token = p.advanceToken()
	varNode, err := p.ParseVar()
	if err != nil {
		return functionNode, fmt.Errorf("function parser -> %s", err)
	}
	functionNode.input = varNode

	token = p.advanceToken()
	if token.Type != PERIOD_TOKEN {
		return functionNode, fmt.Errorf("function parser -> excpected PERIOD, got %s at %d", token.Type, token.Position)
	}

	token = p.advanceToken()
	expression, err := p.ParseExp(0)
	if err != nil {
		return functionNode, fmt.Errorf("function parser -> %s", err)
	}
	functionNode.body = expression

	return functionNode, nil
}

func (p *EngParser) ParseGrouped() (Node, error) {

	token := p.currentToken()
	if token.Type != LPAREN_TOKEN {
		return ProgramNode{}, fmt.Errorf("group parser -> excpected LPAREN, got %s at %d", token.Type, token.Position)
	}

	token = p.advanceToken()
	node, err := p.ParseExp(0)
	if err != nil {
		return node, fmt.Errorf("group parser -> %s", err)
	}

	token = p.advanceToken()
	if token.Type != RPAREN_TOKEN {
		return node, fmt.Errorf("group parser -> excpected RPAREN, got %s at %d", token.Type, token.Position)
	}

	return node, nil
}
