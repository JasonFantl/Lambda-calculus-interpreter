package interpreter

import "fmt"

type Evaluator struct {
	storedFuncs map[string]Node
	showSteps   bool
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		storedFuncs: make(map[string]Node),
		showSteps:   false,
	}
}

func (e *Evaluator) Evaluate(node Node, showSteps bool) Node {
	e.showSteps = showSteps
	if showSteps {
		fmt.Printf("Eval sequence bottom up:\n")
	}
	eval := e.evalNode(node, make(map[string]Node), 1)

	if eval != nil {
		// check if equal to any stored
		for name, node := range e.storedFuncs {
			if equalNodes(eval, node) {
				return NameNode{name}
			}
		}
	}

	return eval

}

func (e *Evaluator) evalNode(node Node, bindings map[string]Node, p int) Node {
	if p > 999 {
		return ErrorNode{Err: fmt.Errorf("to many evaluations, terminating")}
	}
	var eval Node
	switch node := node.(type) {
	case VarNode:
		eval = e.evalVar(node, bindings, p)
	case NameNode:
		eval = e.evalName(node, bindings, p)
	case NamedFunctionNode:
		eval = e.evalNamedFunction(node, bindings, p)
	case FunctionNode:
		eval = e.evalFunction(node, bindings, p)
	case ApplicationNode:
		eval = e.evalApplication(node, bindings, p)
	}

	return eval
}

func (e *Evaluator) evalVar(varNode VarNode, bindings map[string]Node, p int) Node {
	e.Log(p, fmt.Sprintf("%s =>", varNode))

	eval, exists := bindings[varNode.identifier]

	// check if it needs to be marked
	if !exists {
		if varNode.identifier[0] != '*' {
			markedName := "*" + varNode.identifier
			e.Log(p, fmt.Sprintf("Warning: var %s is not bound, renaming to %s to mark as free", varNode, markedName))
			varNode.identifier = markedName
		}
		eval = varNode // var can evaluate to itself if not bound
	}

	e.Log(p, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) evalName(nameNode NameNode, bindings map[string]Node, p int) Node {
	e.Log(p, fmt.Sprintf("%s =>", nameNode))
	var eval Node

	eval, exists := e.storedFuncs[nameNode.identifier]
	if !exists {
		fmt.Printf("name %s is not defined", nameNode.identifier)
	}

	e.Log(p, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) evalNamedFunction(namedFunctionNode NamedFunctionNode, bindings map[string]Node, p int) Node {
	e.Log(p, fmt.Sprintf("%s =>", namedFunctionNode))
	var eval Node

	// check name not already taken
	s := namedFunctionNode.name.identifier
	_, exists := e.storedFuncs[s]
	if exists {
		e.Log(p, fmt.Sprintf("name %s is already defined", s))
	} else {
		betaReduced := e.evalNode(namedFunctionNode.body, bindings, p+1)
		e.storedFuncs[s] = betaReduced
		eval = betaReduced
	}

	e.Log(p, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) evalFunction(functionNode FunctionNode, bindings map[string]Node, p int) Node {
	e.Log(p, fmt.Sprintf("%s =>", functionNode))
	var eval Node

	// check if input name is already bound
	s := functionNode.input.identifier
	pre, exist := bindings[s]
	if exist {
		e.Log(p, fmt.Sprintf("Warning: function var %s is already defined, locally overridding", s))
	}
	e.Log(p, fmt.Sprintf("Binding %s", s))
	bindings[s] = functionNode.input

	// eval
	functionNode.body = e.evalNode(functionNode.body, bindings, p+1)
	eval = functionNode

	// unbind
	e.Log(p, fmt.Sprintf("Unbinding %s", s))
	if exist { // to remove local scope binding
		bindings[s] = pre
	} else {
		delete(bindings, s)
	}

	e.Log(p, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) evalApplication(applicationNode ApplicationNode, bindings map[string]Node, p int) Node {
	e.Log(p, fmt.Sprintf("%s =>", applicationNode))
	var eval Node

	applicationNode.lExp = e.evalNode(applicationNode.lExp, bindings, p+1)
	e.Log(p, fmt.Sprintf("=> %s =>", applicationNode))

	applicationNode.rExp = e.evalNode(applicationNode.rExp, bindings, p+1)
	e.Log(p, fmt.Sprintf("=> %s =>", applicationNode))

	switch f := applicationNode.lExp.(type) {
	case FunctionNode:
		s := f.input.identifier

		// bind variable
		pre, exist := bindings[s]
		if exist {
			e.Log(p, fmt.Sprintf("Warning: var %s is already defined, locally overridding", s))
		}
		bindings[s] = applicationNode.rExp
		e.Log(p, fmt.Sprintf("Binding %s to %s", s, bindings[s]))

		// eval
		eval = e.evalNode(f.body, bindings, p+1)

		// unbind
		e.Log(p, fmt.Sprintf("Unbinding %s from %s", s, bindings[s]))
		if exist { // to remove local scope binding
			bindings[s] = pre
		} else {
			delete(bindings, s)
		}
	default:
		e.Log(p, fmt.Sprintf("application %s cannot be simplified", applicationNode))
		eval = applicationNode
	}

	e.Log(p, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) Log(n int, s string) {
	if e.showSteps {
		for i := 0; i < n; i++ {
			fmt.Print("\u2502 ")
		}
		fmt.Println(s)
	}
}
