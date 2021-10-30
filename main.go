package main

import (
	"fmt"
	"os"
	"teenytinycompiler/lexer"
	"teenytinycompiler/parser"
)

func main() {
	fmt.Println("Teeny Tiny Compiler")

	if len(os.Args) < 2 {
		fmt.Println("Error: Compiler needs source file as argument")
		os.Exit(1)
	}

	input, err := os.ReadFile(os.Args[1])
	check(err)

	lex := lexer.Constructor(string(input))
	parse := parser.Constructor(lex)
	parse.Program()

	fmt.Println("Parsing completed.")
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
