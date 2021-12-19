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

type NamedFunctionDefNode struct {
	identifier string
	inputs     []VarNode
	body       Node
}

func (n NamedFunctionDefNode) String() string {
	inputs := ""
	for _, input := range n.inputs {
		inputs += input.String() + " "
	}

	return fmt.Sprintf("%s %s= %s", n.identifier, inputs, n.body)
}

type NamedFunctionCallNode struct {
	identifier string
	inputs     []Node
}

func (n NamedFunctionCallNode) String() string {
	inputs := ""
	for _, input := range n.inputs {
		inputs += "(" + input.String() + ") "
	}

	return fmt.Sprintf("( %s %s)", n.identifier, inputs)
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
