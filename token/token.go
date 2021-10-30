package token

type Token struct {
	Text []rune
	Kind Type
}

type Type int

const (
	EOF Type = iota
	NEWLINE
	NUMBER
	IDENT
	STRING

	LABEL
	GOTO
	PRINT
	INPUT
	LET
	IF
	THEN
	ENDIF
	WHILE
	REPEAT
	ENDWHILE

	EQ
	PLUS
	MINUS
	ASTERISK
	SLASH
	EQEQ
	NOTEQ
	LT
	LTEQ
	GT
	GTEQ
)

func (t Type) String() string {
	return [...]string{"EOF", "NEWLINE", "NUMBER", "IDENT", "STRING",
		"LABEL", "GOTO", "PRINT", "INPUT", "LET", "IF", "THEN", "ENDIF", "WHILE", "REPEAT", "ENDWHILE",
		"EQ", "PLUS", "MINUS", "ASTERISK", "SLASH", "EQEQ", "NOTEQ", "LT", "LTEQ", "GT", "GTEQ"}[t]
}

func CheckIfKeyword(keyword string) Type {
	switch keyword {
	case "LABEL":
		return LABEL
	case "GOTO":
		return GOTO
	case "PRINT":
		return PRINT
	case "INPUT":
		return INPUT
	case "LET":
		return LET
	case "IF":
		return IF
	case "THEN":
		return THEN
	case "ENDIF":
		return ENDIF
	case "WHILE":
		return WHILE
	case "REPEAT":
		return REPEAT
	case "ENDWHILE":
		return ENDWHILE
	}
	return -1
}
