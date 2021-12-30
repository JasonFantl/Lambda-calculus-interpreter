package interpreter

import (
	"strconv"
)

// not efficient, but works. Assumes nodes are beta reduced
func equalNodes(n1, n2 Node) bool {
	dBNode1 := toDeBruijnIndices(n1, make(map[string]int), 0)
	dBNode2 := toDeBruijnIndices(n2, make(map[string]int), 0)

	return dBNode1.String() == dBNode2.String()
}

func toDeBruijnIndices(node Node, bindings map[string]int, level int) Node {
	var eval Node
	switch node := node.(type) {
	case VarNode:
		eval = dBVar(node, bindings, level)
	case FunctionNode:
		eval = dBFunction(node, bindings, level)
	case ApplicationNode:
		eval = dBApplication(node, bindings, level)
	}

	return eval
}

func dBFunction(functionNode FunctionNode, bindings map[string]int, level int) Node {
	var eval Node

	// check if input name is already bound
	s := functionNode.input.identifier
	pre, exist := bindings[s]

	bindings[s] = -level

	// eval
	functionNode.input.identifier = ""
	functionNode.body = toDeBruijnIndices(functionNode.body, bindings, level+1)
	eval = functionNode

	// unbind
	if exist { // to remove local scope binding
		bindings[s] = pre
	} else {
		delete(bindings, s)
	}

	return eval
}

func dBApplication(applicationNode ApplicationNode, bindings map[string]int, level int) Node {
	return ApplicationNode{toDeBruijnIndices(applicationNode.lExp, bindings, level), toDeBruijnIndices(applicationNode.rExp, bindings, level)}
}

func dBVar(varNode VarNode, bindings map[string]int, level int) Node {

	l, exists := bindings[varNode.identifier]

	// check if it needs to be marked
	if !exists {
		// fmt.Printf("Warning: equality checker found an unbound var %s\n", varNode)
	} else {
		varNode.identifier = strconv.Itoa(l + level)
	}

	return varNode
}
