package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jasonfantl/lambda_interpreter/interpreter"
)

func main() {

	// check if interpreter or read file
	filePtr := flag.String("file", "/", "filepath to file to run")
	flag.Parse()

	// default is shell
	reader := bufio.NewReader(os.Stdin)

	if *filePtr != "/" {
		f, err := os.Open(*filePtr)
		if err != nil {
			fmt.Println(err)
			return
		}

		reader = bufio.NewReader(f)
	}

	// REPL
	running := true
	for running {
		// only print input if using shell
		if *filePtr == "/" {
			fmt.Print(">>")
		}

		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if err == io.EOF {
			running = false
		}

		if input == "exit\n" {
			break
		}

		l := interpreter.NewLexer(input)
		tokens, err := l.Tokenize()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tokens)

		p := interpreter.NewEngParser(tokens)
		tree, err := p.Parse()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tree)
	}
}
