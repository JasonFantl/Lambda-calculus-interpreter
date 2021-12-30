package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jasonfantl/lambda_interpreter/interpreter"
)

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
		if input == "exit\n" {
			break
		}
		if input[0] == ':' {
			if input[0:2] == ":l" {
				loadFile(input[3 : len(input)-1])
			}
			continue
		}

		fmt.Println(interpret(input))
	}
}

func loadFile(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(f)

	for {
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}

		interpret(input)
	}

	fmt.Printf("Loaded file: %s\n", filepath)

}

func interpret(input string) interpreter.Node {
	l := interpreter.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(tokens)

	p := interpreter.NewParser(tokens)
	tree, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(tree)

	return evaluator.Evaluate(tree)
}
