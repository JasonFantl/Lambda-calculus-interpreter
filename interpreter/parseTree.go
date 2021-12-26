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

type NameNode struct {
	identifier string
}

func (n NameNode) String() string {
	return string(n.identifier)
}

type NamedFunctionNode struct {
	identifier NameNode
	function   FunctionNode
}

func (n NamedFunctionNode) String() string {
	return fmt.Sprintf("%s = %s", n.identifier, n.function)
}

type ApplicationNode struct {
	lExp Node
	rExp Node
}

func (n ApplicationNode) String() string {
	return fmt.Sprintf("( %s %s )", n.lExp, n.rExp)
}

type FunctionNode struct {
	inputs []VarNode
	body   Node
}

func (n FunctionNode) String() string {
	inputs := ""
	for _, node := range n.inputs {
		inputs += node.String() + " "
	}
	return fmt.Sprintf("\\ %s . %s", inputs, n.body)
}
