package main

import (
	"fmt"
	"interpreter/lexer"
	"os"
)

func main() {
	f, err := os.ReadFile("test.txt")
	if err != nil {
		panic("could not open file")
	}
	l := lexer.New(string(f))
	tokens := l.Tokenize()
	for _, v := range tokens {
		fmt.Println(v)
	}
}
