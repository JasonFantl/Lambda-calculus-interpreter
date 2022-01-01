package interpreter

import (
	"fmt"
)

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

	if node == nil {
		return ErrorNode{fmt.Errorf("")}
	}

	switch node := node.(type) {
	case NamedFunctionNode: // if its a function definition, store
		s := node.name.identifier
		_, exists := e.storedFuncs[s]
		if exists {
			return ErrorNode{fmt.Errorf("name %s is already defined", s)}
		}

		eval := e.evalNode(e.toDeBruijn(node.body), make([]Node, 0), BoundVarDB{0}, 0, 0)
		e.storedFuncs[s] = eval
		return eval
	}

	// otherwise evaluate
	eval := e.evalNode(e.toDeBruijn(node), make([]Node, 0), BoundVarDB{0}, 0, 0)
	// check if result is recognized
	for name, node := range e.storedFuncs {
		if node.String() == eval.String() {
			return NameNode{name}
		}
	}
	return eval
}

func (e *Evaluator) evalNode(node Node, bindings []Node, en Node, m, p int) Node {
	if p > 999 {
		return ErrorNode{Err: fmt.Errorf("stack overflow")}
	}

	e.Log(p, fmt.Sprintf("%s =>", node))

	var eval Node
	switch node := node.(type) {
	case BoundVarDB:
		if node.index == m {
			eval = e.evalNode(en, bindings, en, m, p)
		} else {
			eval = node
		}
	case FreeVarDB:
		eval = node
	case FunctionDB:
		eval = e.evalFunction(node, bindings, e.shiftIndecies(en, 1, 0), m+1, p)
	case ApplicationDB:
		eval = e.evalApplication(node, bindings, en, m, p)
	case ErrorNode:
		eval = node
	default:
		return ErrorNode{fmt.Errorf("unrecognized node type %s", node.String())}
	}

	e.Log(p, fmt.Sprintf("=> %s", eval))
	return eval
}

func (e *Evaluator) evalFunction(functionNode FunctionDB, bindings []Node, en Node, m, p int) Node {

	newBindings := append(bindings, VarNode{"~"})
	e.Log(p, fmt.Sprintf("binding ~ %v", newBindings))

	functionNode.body = e.evalNode(functionNode.body, newBindings, en, m, p+1) // special varNode to indicate self
	eval := functionNode

	e.Log(p, fmt.Sprintf("unbinding ~ %v", bindings))

	return eval
}

func (e *Evaluator) evalApplication(applicationNode ApplicationDB, bindings []Node, en Node, m, p int) Node {

	e.Log(p, fmt.Sprintf("evaluating the left expression"))
	applicationNode.lExp = e.evalNode(applicationNode.lExp, bindings, en, m, p+1)
	e.Log(p, fmt.Sprintf("=> %s =>", applicationNode))

	e.Log(p, fmt.Sprintf("evaluating the right expression"))
	applicationNode.rExp = e.evalNode(applicationNode.rExp, bindings, en, m, p+1)
	e.Log(p, fmt.Sprintf("=> %s =>", applicationNode))

	switch f := applicationNode.lExp.(type) {
	case FunctionDB:
		e.Log(p, fmt.Sprintf("evaluating the entire expression"))
		newBindings := append(bindings, applicationNode.rExp)
		e.Log(p, fmt.Sprintf("binding %s %v", applicationNode.rExp, newBindings))
		eval := e.evalNode(f.body, newBindings, e.shiftIndecies(applicationNode.rExp, 1, 0), 0, p+1)
		eval = e.shiftIndecies(eval, -1, 0)
		e.Log(p, fmt.Sprintf("unbinding %s %v", applicationNode.rExp, bindings))
		return eval
	}

	return applicationNode
}

func (e *Evaluator) shiftIndecies(node Node, i, c int) Node {
	switch node := node.(type) {
	case FreeVarDB:
		return node
	case BoundVarDB:
		if node.index < c {
			return node
		}
		node.index += i
		return node
	case FunctionDB:
		node.body = e.shiftIndecies(node.body, i, c+1)
		return node
	case ApplicationDB:
		node.lExp = e.shiftIndecies(node.lExp, i, c)
		node.rExp = e.shiftIndecies(node.rExp, i, c)
		return node
	case ErrorNode:
		return node
	}
	return ErrorNode{fmt.Errorf("unrecognized type for shifting %s", node)}
}

func (e *Evaluator) Log(n int, s string) {
	if e.showSteps {
		for i := 0; i < n; i++ {
			fmt.Print("\u2502 ")
		}
		fmt.Println(s)
	}
}

func (e *Evaluator) equalNodes(n1, n2 Node) bool {
	dBNode1 := e.toDeBruijn(n1)
	dBNode2 := e.toDeBruijn(n2)

	return dBNode1.String() == dBNode2.String()
}
