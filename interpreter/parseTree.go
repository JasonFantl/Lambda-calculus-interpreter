package interpreter

import "fmt"

type Node interface {
	String() string
}

type ErrorNode struct {
	Err error
}

func (n ErrorNode) String() string {
	return n.Err.Error()
}

type VarNode struct {
	identifier string
}

func (n VarNode) String() string {
	return string(n.identifier)
}

type NameNode struct {
	identifier string
}

func (n NameNode) String() string {
	return string(n.identifier)
}

type NamedFunctionNode struct {
	name NameNode
	body Node
}

func (n NamedFunctionNode) String() string {
	return fmt.Sprintf("%s = %s", n.name, n.body)
}

type ApplicationNode struct {
	lExp Node
	rExp Node
}

func (n ApplicationNode) String() string {
	return fmt.Sprintf("( %s <- %s )", n.lExp, n.rExp)
}

type FunctionNode struct {
	input VarNode
	body  Node
}

func (n FunctionNode) String() string {
	return fmt.Sprintf("\\ %s . %s", n.input, n.body)
}
