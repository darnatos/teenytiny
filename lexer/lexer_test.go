package lexer

import (
	"reflect"
	"teenytinycompiler/token"
	"testing"
)

func Test_lexerImpl_getTokens(t *testing.T) {
	type fields struct {
		source  []rune
		curChar rune
		curPos  int
	}
	tests := []struct {
		name   string
		fields fields
		want   []*token.Token
	}{
		{
			name: "test case 1",
			fields: fields{
				source: []rune(`IF+-123 foo*THEN/`),
				curPos: -1,
			},
			want: []*token.Token{
				{Kind: token.IF, Text: nil},
				{Kind: token.PLUS, Text: []rune("+")},
				{Kind: token.MINUS, Text: []rune("-")},
				{Kind: token.NUMBER, Text: []rune("123")},
				{Kind: token.IDENT, Text: []rune("foo")},
				{Kind: token.ASTERISK, Text: []rune("*")},
				{Kind: token.THEN, Text: nil},
				{Kind: token.SLASH, Text: []rune("/")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &lexerImpl{
				source:  tt.fields.source,
				curChar: tt.fields.curChar,
				curPos:  tt.fields.curPos,
			}
			lexer.nextChar()
			if got := lexer.getTokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
