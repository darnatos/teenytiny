package main

import (
	"fmt"
	"os"
	"teenytinycompiler/emitter"
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

	// Initialize the lexer, emitter, and parser
	lex := lexer.Constructor(string(input))
	emit := emitter.Constructor("out.c")
	parse := parser.Constructor(lex, emit)
	parse.Program()
	emit.WriteFile()

	fmt.Println("Parsing completed.")
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
