package lexer

import (
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
	byte c = 0;
	for c < a {
		c = nil;
	}
	fun d(string a){return a;}
}
	`
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
		{token.BYTE, ""},
		{token.IDENTIFIER, "c"},
		{token.ASSIGN, ""},
		{token.FLOAT, "0"},
		{token.SEMICOLON, ""},
		{token.FOR, ""},
		{token.IDENTIFIER, "c"},
		{token.LT, ""},
		{token.IDENTIFIER, "a"},
		{token.LCURLY, ""},
		{token.IDENTIFIER, "c"},
		{token.ASSIGN, ""},
		{token.NIL, ""},
		{token.SEMICOLON, ""},
		{token.RCURLY, ""},
		{token.FUN, ""},
		{token.IDENTIFIER, "d"},
		{token.LPAREN, ""},
		{token.STRINGT, ""},
		{token.IDENTIFIER, "a"},
		{token.RPAREN, ""},
		{token.LCURLY, ""},
		{token.RETURN, ""},
		{token.IDENTIFIER, "a"},
		{token.SEMICOLON, ""},
		{token.RCURLY, ""},
		{token.RCURLY, ""},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tests) != len(tokens) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(tests), len(tokens))
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

func TestTokenizeInvalidTokens(t *testing.T) {
	input := "шчш {} if else ELSE"
	tests := []struct {
		expectedType  token.TokenType
		expectedValue string
	}{
		{token.ERR, ""}, // each cyrlic letter is two bytes so the lexer will output two errors for each letter
		{token.ERR, ""},
		{token.ERR, ""},
		{token.ERR, ""},
		{token.ERR, ""},
		{token.ERR, ""},
		{token.LCURLY, ""},
		{token.RCURLY, ""},
		{token.IF, ""},
		{token.ELSE, ""},
		{token.IDENTIFIER, "ELSE"},
	}
	l := New(input)
	tokens := l.Tokenize()
	for i, tc := range tests {
		if tokens[i].Type != tc.expectedType {
			t.Fatalf("tests[%d] - TokenType wrong. expected=%q, got=%q", i, tc.expectedType, tokens[i].Type)
		}

		if tokens[i].Value != tc.expectedValue {
			t.Fatalf("tests[%d] - token value wrong. expected=%q, got=%q", i, tc.expectedValue, tokens[i].Value)
		}
	}

}

func TestWhitespaceCharacters(t *testing.T) {
	input := `

	     		    



	
	`
	l := New(input)
	tokens := l.Tokenize()
	if len(tokens) != 1 {
		t.Fatalf("wrong number of tokens. expected=1, got=%d", len(tokens))
	}
	if tokens[0].Type != "EOF" {
		t.Fatalf("TokenType wrong. expected=EOF, got=%q", tokens[0].Type)
	}
}
