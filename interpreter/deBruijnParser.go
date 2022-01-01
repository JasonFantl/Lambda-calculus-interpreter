package interpreter

import (
	"fmt"
)

func (e *Evaluator) toDeBruijn(root Node) Node {
	if e.showSteps {
		fmt.Println("intermidiate representation:")
	}
	return e.dBNode(root, make(map[string]int), 0)
}

func (e *Evaluator) dBNode(node Node, bindings map[string]int, level int) Node {
	var eval Node
	switch node := node.(type) {
	case NameNode:
		eval = e.dBName(node, bindings, level)
	case NamedFunctionNode:
		eval = e.dBNamedFunction(node, bindings, level)
	case VarNode:
		eval = e.dBVar(node, bindings, level)
	case FunctionNode:
		eval = e.dBFunction(node, bindings, level)
	case ApplicationNode:
		eval = e.dBApplication(node, bindings, level)
	}

	return eval
}

func (e *Evaluator) dBFunction(functionNode FunctionNode, bindings map[string]int, level int) Node {
	e.Log(level, fmt.Sprintf("%s => function", functionNode))

	// check if input name is already bound
	s := functionNode.input.identifier
	pre, exist := bindings[s]

	bindings[s] = -level

	// eval
	eval := FunctionDB{e.dBNode(functionNode.body, bindings, level+1)}

	// unbind
	if exist { // to remove local scope binding
		bindings[s] = pre
	} else {
		delete(bindings, s)
	}

	e.Log(level, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) dBApplication(applicationNode ApplicationNode, bindings map[string]int, level int) Node {
	e.Log(level, fmt.Sprintf("%s => application", applicationNode))

	eval := ApplicationDB{
		e.dBNode(applicationNode.lExp, bindings, level),
		e.dBNode(applicationNode.rExp, bindings, level),
	}

	e.Log(level, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) dBVar(varNode VarNode, bindings map[string]int, level int) Node {
	e.Log(level, fmt.Sprintf("%s => var", varNode))
	var eval Node

	l, exists := bindings[varNode.identifier]
	if !exists {
		eval = FreeVarDB{varNode.identifier}
	} else {
		eval = BoundVarDB{l + level - 1}
	}

	e.Log(level, fmt.Sprintf("=> %s", eval))
	return eval
}

// assumes storedFuncs stores funcs with DB indexing
func (e *Evaluator) dBName(nameNode NameNode, bindings map[string]int, level int) Node {
	e.Log(level, fmt.Sprintf("%s => name", nameNode))
	var eval Node

	body, exists := e.storedFuncs[nameNode.identifier]
	if !exists {
		eval = ErrorNode{fmt.Errorf("name %s is not defined", nameNode.identifier)}
	} else {
		eval = body
	}

	e.Log(level, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) dBNamedFunction(namedFunctionNode NamedFunctionNode, bindings map[string]int, level int) Node {
	e.Log(level, fmt.Sprintf("%s =>", namedFunctionNode))

	eval := e.dBNode(namedFunctionNode.body, bindings, level)

	e.Log(level, fmt.Sprintf("=> %s", eval))
	return eval
}

type FunctionDB struct {
	body Node
}

func (n FunctionDB) String() string {
	return fmt.Sprintf("[%s]", n.body)
}

type ApplicationDB struct {
	lExp, rExp Node
}

func (n ApplicationDB) String() string {
	return fmt.Sprintf("(%s %s)", n.lExp, n.rExp)
}

type BoundVarDB struct {
	index int
}

func (n BoundVarDB) String() string {
	return fmt.Sprintf("%d", n.index)
}

type FreeVarDB struct {
	identifier string
}

func (n FreeVarDB) String() string {
	return n.identifier
}
