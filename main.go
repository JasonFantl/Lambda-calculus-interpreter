package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jasonfantl/lambda_interpreter/interpreter"
)

// TODO: fix beta reductions: for example ":s FLIP T" evaluates to T due to variable capture
var evaluator *interpreter.Evaluator
var filepaths []string

func main() {

	// check if interpreter or read file
	flag.Parse()

	evaluator = interpreter.NewEvaluator()

	filepaths = flag.Args()
	loadFiles(filepaths)

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

		execute(input)
	}
}

func execute(input string) {
	if len(input) == 0 {
		return
	}
	// special shell commands
	if input[0] == ':' {
		switch input[1] {
		case 'l': // load file
			if len(input) > 3 {
				filepaths = append(filepaths, strings.Fields(input)...)
				loadFiles(filepaths)
			} else {
				fmt.Println("at least one filepath must be given")
			}
		case 'r': // reload file
			loadFiles(filepaths)
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

func loadFiles(filepaths []string) {
	for _, filepath := range filepaths {
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		reader := bufio.NewReader(f)
		fmt.Printf("---------- Loading file: %s... ----------\n", filepath)

		for {
			input, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				break
			}

			fmt.Printf(">>%s", input)
			if err == io.EOF {
				fmt.Println()
			}
			execute(input)

			if err == io.EOF {
				break
			}
		}

		fmt.Printf("---------- Loaded file: %s ----------\n", filepath)
	}
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
