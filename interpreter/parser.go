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

// NOTE: each parse ends at the last sucsessful token of the parse

func (p *Parser) ParseProgram() (ProgramNode, error) {
	program := ProgramNode{make([]Node, 0)}

	for token := p.currentToken(); token.Type != EOF_TOKEN; token = p.advanceToken() {

		// ignore all NEWLINEs before an expression
		if token.Type == NEWLINE_TOKEN {
			continue
		}

		expression, err := p.ParseExpression()
		if err != nil {
			return program, fmt.Errorf("program parser -> %s", err)
		}
		program.nodes = append(program.nodes, expression)
	}

	return program, nil
}

// for all the love that is mighty, please refactor this
func (p *Parser) ParseExpression() (Node, error) {
	token := p.currentToken()

	// Production: <var> S'
	if token.Type == VAR_TOKEN {
		varNode, err := p.ParseVar()
		if err != nil {
			return varNode, fmt.Errorf("expression parser -> %s", err)
		}
		p.advanceToken()
		appNode, err := p.ParseExpressionPrime(varNode)
		if err != nil {
			return appNode, fmt.Errorf("expression parser -> %s", err)
		}
		return appNode, nil
	}

	// Production: \<var> . S S'
	if token.Type == LAMBDA_TOKEN {
		funcNode, err := p.ParseFunction()
		if err != nil {
			return funcNode, fmt.Errorf("expression parser -> %s", err)
		}
		p.advanceToken()
		appNode, err := p.ParseExpressionPrime(funcNode)
		if err != nil {
			return appNode, fmt.Errorf("expression parser -> %s", err)
		}
		return appNode, nil
	}

	// Production: ( S ) S'
	if token.Type == LPAREN_TOKEN {
		expNode, err := p.ParseGrouped()
		if err != nil {
			return expNode, fmt.Errorf("expression parser -> %s", err)
		}
		p.advanceToken()
		appNode, err := p.ParseExpressionPrime(expNode)
		if err != nil {
			return appNode, fmt.Errorf("expression parser -> %s", err)
		}
		return appNode, nil
	}

	return ProgramNode{}, fmt.Errorf("expression parser -> expected VAR, LAMBDA, or LPAREN, got %s at %d", token.Type, token.Position)
}

// we are hardcoding in FIRST, should be automated to be safe
var parseExpressionPrimeFIRST = map[TokenType]bool{VAR_TOKEN: true, LAMBDA_TOKEN: true, LPAREN_TOKEN: true}

func (p *Parser) ParseExpressionPrime(node Node) (Node, error) {

	token := p.currentToken()

	if parseExpressionPrimeFIRST[token.Type] { // this is the application production
		expr, err := p.ParseExpression()
		if err != nil {
			return node, fmt.Errorf("expressionPrime parser -> %s", err)
		}
		return ApplicationNode{node, expr}, nil
	}

	// this is the epsilon production, need to decrement token counter
	p.position--

	return node, nil
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
	expression, err := p.ParseExpression()
	if err != nil {
		return functionNode, fmt.Errorf("function parser -> %s", err)
	}
	functionNode.body = expression

	return functionNode, nil
}

func (p *Parser) ParseGrouped() (Node, error) {

	token := p.currentToken()
	if token.Type != LPAREN_TOKEN {
		return ProgramNode{}, fmt.Errorf("group parser -> excpected LPAREN, got %s at %d", token.Type, token.Position)
	}

	token = p.advanceToken()
	node, err := p.ParseExpression()
	if err != nil {
		return node, fmt.Errorf("group parser -> %s", err)
	}

	token = p.advanceToken()
	if token.Type != RPAREN_TOKEN {
		return node, fmt.Errorf("group parser -> excpected RPAREN, got %s at %d", token.Type, token.Position)
	}

	return node, nil
}
