package interpreter

import "fmt"

type Node interface {
	String() string
}

type ProgramNode struct {
	nodes []Node
}

func (n ProgramNode) String() string {
	result := ""
	for _, node := range n.nodes {
		result += node.String() + "\n"
	}
	return result
}

type VarNode struct {
	identifier string
}

func (n VarNode) String() string {
	return string(n.identifier)
}

type ApplicationNode struct {
	lExp Node
	rExp Node
}

func (n ApplicationNode) String() string {
	return fmt.Sprintf("( %s %s )", n.lExp, n.rExp)
}

type FunctionNode struct {
	input VarNode
	body  Node
}

func (n FunctionNode) String() string {
	return fmt.Sprintf("\\ %s . %s", n.input, n.body)
}
