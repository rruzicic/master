package main

import (
	"bufio"
	"fmt"
	"interpreter/lexer"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		file(os.Args[1])
	} else {
		panic("wrong number of args")
	}
}

func file(filename string) {
	f, err := os.ReadFile(filename)
	if err != nil {
		panic("could not open file")
	}
	run(string(f))
}

func repl() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("could not read input: %s", err))
		}
		run(input)
	}
}

func run(input string) {
	l := lexer.New(input)
	tokens := l.Tokenize()
	for _, v := range tokens {
		fmt.Println(v)
	}
}
