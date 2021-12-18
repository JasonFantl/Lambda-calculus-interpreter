package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jasonfantl/lambda_interpreter/interpreter"
)

func main() {

	// REPL
	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(">>")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if input == "exit\n" {
			break
		}

		l := interpreter.NewLexer(input)
		tokens, err := l.Tokenize()
		if err != nil {
			fmt.Println(err)
		}

		p := interpreter.NewParser(tokens)
		tree, err := p.Parse()

		if err != nil {
			fmt.Println(err)
			fmt.Println(tree)
		} else {
			fmt.Println(tokens)
			fmt.Println(tree)
		}
	}
}
