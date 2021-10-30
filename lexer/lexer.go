package lexer

import (
	"fmt"
	"os"
	"teenytinycompiler/token"
	"unicode"
)

type Lexer interface {
	GetToken() *token.Token
}

type lexerImpl struct {
	source  []rune
	curChar rune
	curPos  int
}

func Constructor(str string) Lexer {
	lex := lexerImpl{source: []rune(str), curPos: -1}
	lex.nextChar()
	return &lex
}

func (lexer *lexerImpl) nextChar() {
	lexer.curPos++
	if lexer.curPos >= len(lexer.source) {
		lexer.curChar = rune(0)
	} else {
		lexer.curChar = lexer.source[lexer.curPos]
	}
}

func (lexer lexerImpl) peek() rune {
	if lexer.curPos+1 >= len(lexer.source) {
		return rune(0)
	}
	return lexer.source[lexer.curPos+1]
}

func (lexer lexerImpl) abort(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func (lexer *lexerImpl) skipWhitespace() {
	for lexer.curChar == '\t' || lexer.curChar == '\r' || lexer.curChar == ' ' {
		lexer.nextChar()
	}
}

func (lexer *lexerImpl) skipComment() {
	if lexer.curChar == '#' {
		for lexer.curChar != '\n' && lexer.curChar != rune(0) {
			lexer.nextChar()
		}
	}
}

func (lexer *lexerImpl) GetToken() *token.Token {
	lexer.skipComment()
	lexer.skipWhitespace()

	tok := &token.Token{Text: []rune{lexer.curChar}}
	switch lexer.curChar {
	case '+':
		tok.Kind = token.PLUS
	case '-':
		tok.Kind = token.MINUS
	case '*':
		tok.Kind = token.ASTERISK
	case '/':
		tok.Kind = token.SLASH
	case '\n':
		tok.Kind = token.NEWLINE
	case rune(0):
		tok.Text = nil
		tok.Kind = token.EOF
	case '=':
		if lexer.peek() == '=' {
			tok.Kind = token.EQEQ
			lexer.nextChar()
			tok.Text = append(tok.Text, lexer.curChar)
		} else {
			tok.Kind = token.EQ
		}
	case '<':
		if lexer.peek() == '=' {
			tok.Kind = token.LTEQ
			lexer.nextChar()
			tok.Text = append(tok.Text, lexer.curChar)
		} else {
			tok.Kind = token.LT
		}
	case '>':
		if lexer.peek() == '=' {
			tok.Kind = token.GTEQ
			lexer.nextChar()
			tok.Text = append(tok.Text, lexer.curChar)
		} else {
			tok.Kind = token.GT
		}
	case '!':
		if lexer.peek() == '=' {
			tok.Kind = token.NOTEQ
			lexer.nextChar()
			tok.Text = append(tok.Text, lexer.curChar)
		} else {
			lexer.abort("Expected !=, got !" + string(lexer.peek()))
		}
	case '"':
		lexer.nextChar()
		startPos := lexer.curPos
		for lexer.curChar != '"' {
			switch lexer.curChar {
			case '\t', '\r', '\n', '%', rune(0):
				lexer.abort("Illegal character in string.")
			default:
				lexer.nextChar()
			}
		}
		tok.Kind = token.STRING
		tok.Text = lexer.source[startPos:lexer.curPos]
	default:
		if unicode.IsDigit(lexer.curChar) {
			startPos := lexer.curPos
			for unicode.IsDigit(lexer.peek()) {
				lexer.nextChar()
			}
			if lexer.peek() == '.' {
				lexer.nextChar()
				if !unicode.IsDigit(lexer.peek()) {
					lexer.abort("Illegal character in number")
				}
				for unicode.IsDigit(lexer.peek()) {
					lexer.nextChar()
				}
			}
			tok.Kind = token.NUMBER
			tok.Text = lexer.source[startPos : lexer.curPos+1]
		} else if unicode.IsLetter(lexer.curChar) {
			startPos := lexer.curPos
			for unicode.IsLetter(lexer.peek()) {
				lexer.nextChar()
			}
			tokText := lexer.source[startPos : lexer.curPos+1]
			keyword := token.CheckIfKeyword(string(tokText))
			if keyword < 0 {
				tok.Kind = token.IDENT
				tok.Text = tokText
			} else {
				tok.Kind = keyword
				tok.Text = nil
			}
		} else {
			lexer.abort("Unknown token: " + string(lexer.curChar))
		}
	}

	lexer.nextChar()
	return tok
}

func (lexer *lexerImpl) getTokens() []*token.Token {
	toks := make([]*token.Token, 0)
	for {
		tok := lexer.GetToken()
		if tok.Kind == token.EOF {
			break
		}
		toks = append(toks, tok)
	}
	return toks
}
