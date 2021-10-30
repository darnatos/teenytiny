package parser

import (
	"fmt"
	"os"
	"teenytinycompiler/lexer"
	"teenytinycompiler/token"
)

type Parser interface {
	Program()
}

type parserImpl struct {
	lex            lexer.Lexer
	symbols        map[string]struct{}
	labelsDeclared map[string]struct{}
	labelsGotoed   map[string]struct{}
	curToken       *token.Token
	peekToken      *token.Token
}

func Constructor(lex lexer.Lexer) Parser {
	pars := parserImpl{
		lex:            lex,
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
		fmt.Println("STATEMENT-PRINT")
		pars.nextToken()
		if pars.checkToken(token.STRING) {
			pars.nextToken()
		} else {
			pars.expression()
		}
	} else if pars.checkToken(token.IF) {
		fmt.Println("STATEMENT-IF")
		pars.nextToken()
		pars.comparison()

		pars.match(token.THEN)
		pars.nl()

		for !pars.checkToken(token.ENDIF) {
			pars.statement()
		}
		pars.match(token.ENDIF)
	} else if pars.checkToken(token.WHILE) {
		fmt.Println("STATEMENT-WHILE")
		pars.nextToken()
		pars.comparison()

		pars.match(token.REPEAT)
		pars.nl()

		for !pars.checkToken(token.ENDWHILE) {
			pars.statement()
		}
		pars.match(token.ENDWHILE)
	} else if pars.checkToken(token.LABEL) {
		fmt.Println("STATEMENT-LABEL")
		pars.nextToken()

		text := string(pars.curToken.Text)
		if _, ok := pars.labelsDeclared[text]; ok {
			pars.abort("Label already exists: " + text)
		}
		pars.labelsDeclared[text] = struct{}{}

		pars.match(token.IDENT)
	} else if pars.checkToken(token.GOTO) {
		fmt.Println("STATEMENT-GOTO")
		pars.nextToken()
		pars.labelsGotoed[string(pars.curToken.Text)] = struct{}{}
		pars.match(token.IDENT)
	} else if pars.checkToken(token.LET) {
		fmt.Println("STATEMENT-LET")
		pars.nextToken()

		text := string(pars.curToken.Text)
		if _, ok := pars.symbols[text]; !ok {
			pars.symbols[text] = struct{}{}
		}

		pars.match(token.IDENT)
		pars.match(token.EQ)

		pars.expression()
	} else if pars.checkToken(token.INPUT) {
		fmt.Println("STATEMENT-INPUT")
		pars.nextToken()

		text := string(pars.curToken.Text)
		if _, ok := pars.symbols[text]; !ok {
			pars.symbols[text] = struct{}{}
		}

		pars.match(token.IDENT)
	} else {
		pars.abort("Invalid statement at " + string(pars.curToken.Text) + " (" + pars.curToken.Kind.String() + ")")
	}
	pars.nl()
}

func (pars *parserImpl) expression() {
	fmt.Println("EXPRESSION")
	pars.term()
	for pars.checkToken(token.PLUS) || pars.checkToken(token.MINUS) {
		pars.nextToken()
		pars.term()
	}
}

func (pars *parserImpl) comparison() {
	fmt.Println("COMPARISON")
	pars.expression()
	if pars.isComparisonOperator() {
		pars.nextToken()
		pars.expression()
	} else {
		pars.abort("Expected comparison operator at: " + string(pars.curToken.Text))
	}

	for pars.isComparisonOperator() {
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

func (pars *parserImpl) term() {
	fmt.Println("TERM")
	pars.unary()
	for pars.checkToken(token.ASTERISK) || pars.checkToken(token.SLASH) {
		pars.nextToken()
		pars.unary()
	}
}

func (pars *parserImpl) unary() {
	fmt.Println("UNARY")
	if pars.checkToken(token.PLUS) || pars.checkToken(token.MINUS) {
		pars.nextToken()
	}
	pars.primary()
}

func (pars *parserImpl) primary() {
	fmt.Println("PRIMARY (" + string(pars.curToken.Text) + ")")
	if pars.checkToken(token.NUMBER) {
		pars.nextToken()
	} else if pars.checkToken(token.IDENT) {
		text := string(pars.curToken.Text)
		if _, ok := pars.symbols[text]; !ok {
			pars.abort("Referencing variable before assignment: " + text)
		}
		pars.nextToken()
	} else {
		pars.abort("Unexpected token at " + string(pars.curToken.Text))
	}
}

func (pars *parserImpl) nl() {
	fmt.Println("NEWLINE")
	pars.match(token.NEWLINE)
	for pars.checkToken(token.NEWLINE) {
		pars.nextToken()
	}
}

func (pars *parserImpl) Program() {
	fmt.Println("PROGRAM")
	for pars.checkToken(token.NEWLINE) {
		pars.nextToken()
	}

	for !pars.checkToken(token.EOF) {
		pars.statement()
	}

	for label := range pars.labelsGotoed {
		if _, ok := pars.labelsDeclared[label]; !ok {
			pars.abort("Attempting to GOTO to undeclared label: " + label)
		}
	}
}
