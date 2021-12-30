package interpreter

import (
	"fmt"
)

type Parser struct {
	tokens   []Token
	position int // token index
	root     Node
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *Parser) advanceToken() Token {
	p.position++
	return p.currentToken()
}

func (p *Parser) currentToken() Token {
	if p.position < len(p.tokens) {
		return p.tokens[p.position]
	}
	return p.tokens[len(p.tokens)-1] // return last token EOF
}

// NOTE: each parse ends at the last successful token of the parse

func (p *Parser) Parse() (Node, error) {

	token := p.currentToken()

	// ignore NEWLINEs before an expression
	for token.Type == NEWLINE_TOKEN {
		token = p.advanceToken()
	}
	if token.Type == DEF_TOKEN {
		return p.ParseNamedFunction()
	}
	if ok := parseExpFIRST[token.Type]; ok {
		return p.ParseExp(0)
	}
	if token.Type == EOF_TOKEN {
		return nil, nil
	}
	return nil, fmt.Errorf("parser -> expected DEF, VAR, LAMBDA, or LPAREN, got %s at %d", token.Type, token.Position)
}

// https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm#climbing

var parseExpFIRST = map[TokenType]bool{
	VAR_TOKEN:    true,
	LAMBDA_TOKEN: true,
	LPAREN_TOKEN: true,
	NAME_TOKEN:   true,
}

func (p *Parser) ParseExp(prec int) (Node, error) {
	exp, err := p.ParsePrefix()
	if err != nil {
		return nil, fmt.Errorf("expression parser -> %s", err)
	}

	// we use FIRST to check if we have an "application operator"
	for nextToken := p.advanceToken(); parseExpFIRST[nextToken.Type] &&
		prec <= 0; nextToken = p.advanceToken() {

		q := prec + 1

		appNode, err := p.ParseExp(q)
		if err != nil {
			return nil, fmt.Errorf("expression parser -> %s", err)
		}

		exp = ApplicationNode{exp, appNode}
	}
	// we went too far, need to decrement
	p.position--

	return exp, nil
}

func (p *Parser) ParsePrefix() (Node, error) {
	token := p.currentToken()

	if token.Type == VAR_TOKEN {
		varNode, err := p.ParseVar()
		if err != nil {
			return varNode, fmt.Errorf("prefix parser -> %s", err)
		}
		return varNode, nil
	}

	if token.Type == LAMBDA_TOKEN {
		funcNode, err := p.ParseFunction()
		if err != nil {
			return funcNode, fmt.Errorf("prefix parser -> %s", err)
		}
		return funcNode, nil
	}

	if token.Type == LPAREN_TOKEN {
		groupedNode, err := p.ParseGrouped()
		if err != nil {
			return groupedNode, fmt.Errorf("prefix parser -> %s", err)
		}
		return groupedNode, nil
	}

	if token.Type == NAME_TOKEN {
		fNameNode, err := p.ParseName()
		if err != nil {
			return fNameNode, fmt.Errorf("prefix parser -> %s", err)
		}
		return fNameNode, nil
	}

	return nil, fmt.Errorf("prefix parser -> expected VAR, LAMBDA, LPAREN, or NAME, got %s at %d", token.Type, token.Position)
}

func (p *Parser) ParseVar() (VarNode, error) {
	varNode := VarNode{}

	token := p.currentToken()
	if token.Type != VAR_TOKEN {
		return varNode, fmt.Errorf("var parser -> excpected VAR, got %s at %d", token.Type, token.Position)
	}
	varNode.identifier = token.Literal

	return varNode, nil
}

func (p *Parser) ParseFunction() (FunctionNode, error) {
	functionNode := FunctionNode{}

	token := p.currentToken()
	if token.Type != LAMBDA_TOKEN {
		return functionNode, fmt.Errorf("function parser -> excpected LAMBDA, got %s at %d", token.Type, token.Position)
	}

	token = p.advanceToken()
	varNodes, err := p.ParseVarList()
	if err != nil {
		return functionNode, fmt.Errorf("function parser -> %s", err)
	}
	// this is weird, but I want to transform these into single var functions, but first need the body

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

	// alright, now we use the inputs to generate nested function nodes
	functionNode.input = varNodes[len(varNodes)-1]
	for i := len(varNodes) - 2; i >= 0; i-- {
		functionNode = FunctionNode{varNodes[i], functionNode}
	}

	return functionNode, nil
}

func (p *Parser) ParseGrouped() (Node, error) {

	token := p.currentToken()
	if token.Type != LPAREN_TOKEN {
		return nil, fmt.Errorf("group parser -> excpected LPAREN, got %s at %d", token.Type, token.Position)
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

func (p *Parser) ParseVarList() ([]VarNode, error) {
	varNodes := make([]VarNode, 0)

	// requires at least one var
	varNode, err := p.ParseVar()
	if err != nil {
		return varNodes, fmt.Errorf("varList parser -> %s", err)
	}
	varNodes = append(varNodes, varNode)

	// the rest are optional. FIRST here is just VAR_TOKEN
	for token := p.advanceToken(); token.Type == VAR_TOKEN; token = p.advanceToken() {
		varNode, err := p.ParseVar()
		if err != nil {
			return varNodes, fmt.Errorf("varList parser -> %s", err)
		}
		varNodes = append(varNodes, varNode)
	}
	p.position-- // we overshot by one

	return varNodes, nil
}

func (p *Parser) ParseName() (NameNode, error) {
	nameNode := NameNode{}

	token := p.currentToken()
	if token.Type != NAME_TOKEN {
		return nameNode, fmt.Errorf("name parser -> excpected NAME, got %s at %d", token.Type, token.Position)
	}
	nameNode.identifier = token.Literal

	return nameNode, nil
}

func (p *Parser) ParseNamedFunction() (NamedFunctionNode, error) {
	namedFunctionNode := NamedFunctionNode{}

	token := p.currentToken()
	if token.Type != DEF_TOKEN {
		return namedFunctionNode, fmt.Errorf("named function parser -> excpected DEF, got %s at %d", token.Type, token.Position)
	}

	token = p.advanceToken()
	nameNode, err := p.ParseName()
	if err != nil {
		return namedFunctionNode, fmt.Errorf("named function parser -> %s", err)
	}
	namedFunctionNode.name = nameNode

	token = p.advanceToken()
	functionNode, err := p.ParseFunction()
	if err != nil {
		return namedFunctionNode, fmt.Errorf("named function parser -> %s", err)
	}
	namedFunctionNode.function = functionNode

	return namedFunctionNode, nil
}
