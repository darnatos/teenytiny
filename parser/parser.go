package parser

import (
	"fmt"
	"os"
	"teenytinycompiler/emitter"
	"teenytinycompiler/lexer"
	"teenytinycompiler/token"
)

type Parser interface {
	Program()
}

type parserImpl struct {
	lex            lexer.Lexer
	emit           emitter.Emitter
	symbols        map[string]struct{}
	labelsDeclared map[string]struct{}
	labelsGotoed   map[string]struct{}
	curToken       *token.Token
	peekToken      *token.Token
}

func Constructor(lex lexer.Lexer, emit emitter.Emitter) Parser {
	pars := parserImpl{
		lex:            lex,
		emit:           emit,
		symbols:        make(map[string]struct{}),
		labelsDeclared: make(map[string]struct{}),
		labelsGotoed:   make(map[string]struct{}),
	}
	pars.nextToken()
	pars.nextToken()
	return &pars
}

func (pars parserImpl) checkToken(kind token.Type) bool {
	return kind == pars.curToken.Kind
}

func (pars parserImpl) checkPeek(kind token.Type) bool {
	return kind == pars.peekToken.Kind
}

func (pars *parserImpl) match(kind token.Type) {
	if pars.checkToken(token.EOF) {
		return
	}
	if !pars.checkToken(kind) {
		pars.abort("Expected: " + kind.String() + ", got: " + pars.curToken.Kind.String())
	}
	pars.nextToken()
}

func (pars *parserImpl) nextToken() {
	pars.curToken, pars.peekToken = pars.peekToken, pars.lex.GetToken()
}

func (pars *parserImpl) abort(message string) {
	fmt.Println("Error! " + message)
	os.Exit(1)
}

func (pars *parserImpl) statement() {
	if pars.checkToken(token.PRINT) {
		// "PRINT" (expression | string)
		pars.nextToken()
		if pars.checkToken(token.STRING) {
			// Simple string, so print it.
			pars.emit.EmitLine("printf(\"" + string(pars.curToken.Text) + "\\n\");")
			pars.nextToken()
		} else {
			// Expected a expression and print the result as a float.
			pars.emit.Emit("printf(\"%.2f\\n\", (float)(")
			pars.expression()
			pars.emit.EmitLine("));")
		}
	} else if pars.checkToken(token.IF) {
		// "IF" comparison "THEN" block "ENDIF"
		pars.nextToken()
		pars.emit.Emit("if(")
		pars.comparison()

		pars.match(token.THEN)
		pars.emit.EmitLine("){")
		pars.nl()

		// Zero or more statements in the body.
		for !pars.checkToken(token.ENDIF) {
			pars.statement()
		}
		pars.match(token.ENDIF)
		pars.emit.EmitLine("}")

	} else if pars.checkToken(token.WHILE) {
		// "WHILE" comparison "REPEAT" block "ENDWHILE"
		pars.nextToken()
		pars.emit.Emit("while(")
		pars.comparison()

		pars.match(token.REPEAT)
		pars.emit.EmitLine("){")
		pars.nl()

		// Zero or more statements in the loop body.
		for !pars.checkToken(token.ENDWHILE) {
			pars.statement()
		}
		pars.match(token.ENDWHILE)
		pars.emit.EmitLine("}")

	} else if pars.checkToken(token.LABEL) {
		// "LABEL" ident
		pars.nextToken()

		text := string(pars.curToken.Text)
		if _, ok := pars.labelsDeclared[text]; ok {
			pars.abort("Label already exists: " + text)
		}
		pars.labelsDeclared[text] = struct{}{}
		pars.emit.EmitLine(string(pars.curToken.Text) + ":")
		pars.match(token.IDENT)

	} else if pars.checkToken(token.GOTO) {
		// "GOTO" ident
		pars.nextToken()
		pars.labelsGotoed[string(pars.curToken.Text)] = struct{}{}
		pars.emit.EmitLine("goto " + string(pars.curToken.Text) + ";")
		pars.match(token.IDENT)

	} else if pars.checkToken(token.LET) {
		// "LET" ident = expression
		pars.nextToken()

		// Check if ident exists in symbol table. If not, declare it.
		text := string(pars.curToken.Text)
		if _, ok := pars.symbols[text]; !ok {
			pars.symbols[text] = struct{}{}
			pars.emit.HeaderLine("float " + string(pars.curToken.Text) + ";")
		}

		pars.emit.Emit(string(pars.curToken.Text) + " = ")
		pars.match(token.IDENT)
		pars.match(token.EQ)

		pars.expression()
		pars.emit.EmitLine(";")

	} else if pars.checkToken(token.INPUT) {
		// "INPUT" ident
		pars.nextToken()

		text := string(pars.curToken.Text)
		if _, ok := pars.symbols[text]; !ok {
			pars.symbols[text] = struct{}{}
			pars.emit.HeaderLine("float " + string(pars.curToken.Text) + ";")
		}
		// Emit scanf but also validate the input. If invalid, set the variable to 0 and clear the input.
		pars.emit.EmitLine("if(0 == scanf(\"%f\", &" + string(pars.curToken.Text) + ")) {")
		pars.emit.EmitLine(string(pars.curToken.Text) + " = 0;")
		pars.emit.EmitLine("scanf(\"%*s\");")
		pars.emit.EmitLine("}")
		pars.match(token.IDENT)

	} else {
		pars.abort("Invalid statement at " + string(pars.curToken.Text) + " (" + pars.curToken.Kind.String() + ")")
	}
	pars.nl()
}

// expression ::= term {( "-" | "+" ) term}
func (pars *parserImpl) expression() {
	pars.term()
	// Can have 0 or more +/- and expressions.
	for pars.checkToken(token.PLUS) || pars.checkToken(token.MINUS) {
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
		pars.term()
	}
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func (pars *parserImpl) comparison() {
	pars.expression()
	// Must be at least one comparison operator and another expression.
	if pars.isComparisonOperator() {
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
		pars.expression()
	} else {
		pars.abort("Expected comparison operator at: " + string(pars.curToken.Text))
	}
	// Can have 0 or more comparison operator and expressions.
	for pars.isComparisonOperator() {
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
		pars.expression()
	}
}

func (pars *parserImpl) isComparisonOperator() bool {
	if pars == nil {
		return false
	}
	return pars.curToken.Kind >= token.EQEQ && pars.curToken.Kind <= token.GTEQ
}

// term ::= unary {( "/" | "*" ) unary}
func (pars *parserImpl) term() {
	pars.unary()
	// Can have 0 or more *// and expressions.
	for pars.checkToken(token.ASTERISK) || pars.checkToken(token.SLASH) {
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
		pars.unary()
	}
}

// unary ::= ["+" | "-"] primary
func (pars *parserImpl) unary() {
	// Optional unary +/-
	if pars.checkToken(token.PLUS) || pars.checkToken(token.MINUS) {
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
	}
	pars.primary()
}

// primary ::= number | ident
func (pars *parserImpl) primary() {
	if pars.checkToken(token.NUMBER) {
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
	} else if pars.checkToken(token.IDENT) {
		// Ensure the variable already exists.
		text := string(pars.curToken.Text)
		if _, ok := pars.symbols[text]; !ok {
			pars.abort("Referencing variable before assignment: " + text)
		}
		pars.emit.Emit(string(pars.curToken.Text))
		pars.nextToken()
	} else {
		// Error!
		pars.abort("Unexpected token at " + string(pars.curToken.Text))
	}
}

func (pars *parserImpl) nl() {
	pars.match(token.NEWLINE)
	for pars.checkToken(token.NEWLINE) {
		pars.nextToken()
	}
}

func (pars *parserImpl) Program() {
	pars.emit.HeaderLine("#include <stdio.h>")
	pars.emit.HeaderLine("int main(void){")

	for pars.checkToken(token.NEWLINE) {
		pars.nextToken()
	}

	for !pars.checkToken(token.EOF) {
		pars.statement()
	}

	pars.emit.EmitLine("return 0;")
	pars.emit.EmitLine("}")

	for label := range pars.labelsGotoed {
		if _, ok := pars.labelsDeclared[label]; !ok {
			pars.abort("Attempting to GOTO to undeclared label: " + label)
		}
	}
}
