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
		{token.TOKEN_GTE, ""},
		{token.TOKEN_EQUAL, ""},
		{token.TOKEN_LTE, ""},
		{token.TOKEN_NOT_EQUAL, ""},
		{token.TOKEN_ASSIGN, ""},
		{token.TOKEN_LCURLY, ""},
		{token.TOKEN_RCURLY, ""},
		{token.TOKEN_LPAREN, ""},
		{token.TOKEN_RPAREN, ""},
		{token.TOKEN_PLUS, ""},
		{token.TOKEN_MINUS, ""},
		{token.TOKEN_MUL, ""},
		{token.TOKEN_DIV, ""},
		{token.TOKEN_BANG, ""},
		{token.TOKEN_ASSIGN, ""},
		{token.NUMBER, "12312312"},
		{token.TOKEN_PLUS, ""},
		{token.TOKEN_MINUS, ""},
		{token.STRING, "adasda"},
		{token.IDENTIFIER, "asasd"},
		{token.TOKEN_LCURLY, ""},
		{token.TOKEN_INT, ""},
		{token.IDENTIFIER, "a"},
		{token.TOKEN_ASSIGN, ""},
		{token.NUMBER, "123"},
		{token.TOKEN_SEMICOLON, ""},
		{token.TOKEN_IF, ""},
		{token.TOKEN_LPAREN, ""},
		{token.IDENTIFIER, "a"},
		{token.TOKEN_EQUAL, ""},
		{token.NUMBER, "123"},
		{token.TOKEN_RPAREN, ""},
		{token.TOKEN_LCURLY, ""},
		{token.TOKEN_BOOL, ""},
		{token.IDENTIFIER, "b"},
		{token.TOKEN_ASSIGN, ""},
		{token.TOKEN_FALSE, ""},
		{token.TOKEN_SEMICOLON, ""},
		{token.TOKEN_RCURLY, ""},
		{token.TOKEN_ELSE, ""},
		{token.TOKEN_IF, ""},
		{token.TOKEN_LPAREN, ""},
		{token.IDENTIFIER, "b"},
		{token.TOKEN_EQUAL, ""},
		{token.TOKEN_TRUE, ""},
		{token.TOKEN_RPAREN, ""},
		{token.TOKEN_LCURLY, ""},
		{token.TOKEN_STRING, ""},
		{token.IDENTIFIER, "a"},
		{token.TOKEN_ASSIGN, ""},
		{token.STRING, "asdasda"},
		{token.TOKEN_SEMICOLON, ""},
		{token.TOKEN_RCURLY, ""},
		{token.TOKEN_BYTE, ""},
		{token.IDENTIFIER, "c"},
		{token.TOKEN_ASSIGN, ""},
		{token.NUMBER, "0"},
		{token.TOKEN_SEMICOLON, ""},
		{token.TOKEN_FOR, ""},
		{token.IDENTIFIER, "c"},
		{token.TOKEN_LT, ""},
		{token.IDENTIFIER, "a"},
		{token.TOKEN_LCURLY, ""},
		{token.IDENTIFIER, "c"},
		{token.TOKEN_ASSIGN, ""},
		{token.TOKEN_NIL, ""},
		{token.TOKEN_SEMICOLON, ""},
		{token.TOKEN_RCURLY, ""},
		{token.TOKEN_FUN, ""},
		{token.IDENTIFIER, "d"},
		{token.TOKEN_LPAREN, ""},
		{token.TOKEN_STRING, ""},
		{token.IDENTIFIER, "a"},
		{token.TOKEN_RPAREN, ""},
		{token.TOKEN_LCURLY, ""},
		{token.TOKEN_RETURN, ""},
		{token.IDENTIFIER, "a"},
		{token.TOKEN_SEMICOLON, ""},
		{token.TOKEN_RCURLY, ""},
		{token.TOKEN_RCURLY, ""},
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
		{token.TOKEN_LCURLY, ""},
		{token.TOKEN_RCURLY, ""},
		{token.TOKEN_IF, ""},
		{token.TOKEN_ELSE, ""},
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
