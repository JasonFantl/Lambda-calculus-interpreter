package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jasonfantl/lambda_interpreter/interpreter"
)

// TODO: fix beta reductions: for example ":s FLIP T" evaluates to T due to variable capture
var evaluator *interpreter.Evaluator

func main() {

	// check if interpreter or read file
	filePtr := flag.String("file", "/", "filepath to file to run")
	flag.Parse()

	evaluator = interpreter.NewEvaluator()

	if *filePtr != "/" {
		loadFile(*filePtr)
	}

	// REPL
	reader := bufio.NewReader(os.Stdin)
	running := true
	for running {
		fmt.Print(">>")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// special shell commands
		if input[0] == ':' {
			switch input[1] {
			case 'l': // load file
				if len(input) > 3 {
					loadFile(input[3 : len(input)-1])
				} else {
					fmt.Println("a filepath must be given")
				}
			case 'e': // exit
				return
			case 's': // show all the steps
				fmt.Println(interpret(input[2:], true))
			default:
				fmt.Println("unrecognized command")
			}
		} else {
			fmt.Println(interpret(input, false))
		}
	}
}

func loadFile(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(f)
	fmt.Printf("Loading file: %s...\n", filepath)

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}

		fmt.Println(interpret(input, false))

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("Loaded file: %s\n", filepath)

}

func interpret(input string, showSteps bool) interpreter.Node {
	l := interpreter.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		return interpreter.ErrorNode{Err: err}
	}
	if showSteps {
		fmt.Println(tokens)
	}

	p := interpreter.NewParser(tokens)
	tree, err := p.Parse()
	if err != nil {
		return interpreter.ErrorNode{Err: err}
	}
	if showSteps {
		fmt.Println(tree)
	}

	return evaluator.Evaluate(tree, showSteps)
}
