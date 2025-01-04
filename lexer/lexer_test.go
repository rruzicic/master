package lexer

import (
	"fmt"
	"interpreter/token"
	"testing"
)

func TestTokenize(t *testing.T) {
	input := `
>= == <= != = { } (

)

+ - * / ! = 12312312 + - "adasda" asasd
{
    int a = 123;
	if (a == 123) {
		bool b = false;
	} else if (b == true) {
		string a = "asdasda";
	}
	
}
	`
	/*
		func a(int pera) {
		int a[];
		for a < 5 {
			a = a + 1;
		}
		return a;
		}
	*/
	tests := []struct {
		expectedType  token.TokenType
		expectedValue string
	}{
		{token.GTE, ""},
		{token.EQUAL, ""},
		{token.LTE, ""},
		{token.NOT_EQUAL, ""},
		{token.ASSIGN, ""},
		{token.LCURLY, ""},
		{token.RCURLY, ""},
		{token.LPAREN, ""},
		{token.RPAREN, ""},
		{token.PLUS, ""},
		{token.MINUS, ""},
		{token.MUL, ""},
		{token.DIV, ""},
		{token.BANG, ""},
		{token.ASSIGN, ""},
		{token.FLOAT, "12312312"},
		{token.PLUS, ""},
		{token.MINUS, ""},
		{token.STRING, "adasda"},
		{token.IDENTIFIER, "asasd"},
		{token.LCURLY, ""},
		{token.INT, ""},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ""},
		{token.FLOAT, "123"},
		{token.SEMICOLON, ""},
		{token.IF, ""},
		{token.LPAREN, ""},
		{token.IDENTIFIER, "a"},
		{token.EQUAL, ""},
		{token.FLOAT, "123"},
		{token.RPAREN, ""},
		{token.LCURLY, ""},
		{token.BOOL, ""},
		{token.IDENTIFIER, "b"},
		{token.ASSIGN, ""},
		{token.FALSE, ""},
		{token.SEMICOLON, ""},
		{token.RCURLY, ""},
		{token.ELSE, ""},
		{token.IF, ""},
		{token.LPAREN, ""},
		{token.IDENTIFIER, "b"},
		{token.EQUAL, ""},
		{token.TRUE, ""},
		{token.RPAREN, ""},
		{token.LCURLY, ""},
		{token.STRINGT, ""},
		{token.IDENTIFIER, "a"},
		{token.ASSIGN, ""},
		{token.STRING, "asdasda"},
		{token.SEMICOLON, ""},
		{token.RCURLY, ""},
		{token.RCURLY, ""},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()
	fmt.Println(tokens)

	if len(tests) != len(tokens) {
		t.Fatalf("wrong number of tokens. expected=%q, got=%q", len(tests), len(tokens))
	}
	for i, tc := range tests {
		if tokens[i].Type != tc.expectedType {
			t.Fatalf("tests[%d] - TokenType wrong. expected=%q, got=%q", i, tc.expectedType, tokens[i].Type)
		}

		if tokens[i].Value != tc.expectedValue {
			t.Fatalf("tests[%d] - token value wrong. expected=%q, got=%q", i, tc.expectedValue, tokens[i].Value)
		}
	}
}
