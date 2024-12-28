package token

import "fmt"

const (
	PLUS      = "+"
	MINUS     = "-"
	MUL       = "*"
	DIV       = "/"
	BANG      = "!"
	LBRACKET  = "["
	RBRACKET  = "]"
	LPAREN    = "("
	RPAREN    = ")"
	LCURLY    = "{"
	RCURLY    = "}"
	SEMICOLON = ";"
	GT        = ">"
	LT        = "<"
	GTE       = ">="
	LTE       = "<="
	ASSIGN    = "="
	EQUAL     = "=="
	NOT_EQUAL = "!="

	FUN    = "fun"
	NIL    = "nil" // maybe
	IF     = "if"
	ELSE   = "else"
	FOR    = "for"
	WHILE  = "while" // maybe
	RETURN = "return"
	AND    = "and"
	OR     = "or"
	TRUE   = "true"
	FALSE  = "false"

	STRINGT    = "STRING_T" //(keyword)
	INT        = "INT"
	BOOL       = "BOOL"
	BYTE       = "BYTE"
	FLOAT      = "FLOAT"
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"

	EOF = "EOF"
)

type Token struct {
	Type  string
	Value string
	Line  int
	Col   int
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %s Value: '%s' Position: %d:%d", t.Type, t.Value, t.Line, t.Col)
}
