package token

import "fmt"

type TokenType string

const (
	TOKEN_PLUS      = "+"
	TOKEN_MINUS     = "-"
	TOKEN_MUL       = "*"
	TOKEN_DIV       = "/"
	TOKEN_BANG      = "!"
	TOKEN_LBRACKET  = "["
	TOKEN_RBRACKET  = "]"
	TOKEN_LPAREN    = "("
	TOKEN_RPAREN    = ")"
	TOKEN_LCURLY    = "{"
	TOKEN_RCURLY    = "}"
	TOKEN_SEMICOLON = ";"
	TOKEN_GT        = ">"
	TOKEN_LT        = "<"
	TOKEN_GTE       = ">="
	TOKEN_LTE       = "<="
	TOKEN_ASSIGN    = "="
	TOKEN_EQUAL     = "=="
	TOKEN_NOT_EQUAL = "!="

	TOKEN_FUN    = "fun"
	TOKEN_NIL    = "nil"
	TOKEN_IF     = "if"
	TOKEN_ELSE   = "else"
	TOKEN_FOR    = "for"
	TOKEN_WHILE  = "while"
	TOKEN_RETURN = "return"
	TOKEN_AND    = "and"
	TOKEN_OR     = "or"
	TOKEN_TRUE   = "true"
	TOKEN_FALSE  = "false"

	TOKEN_STRING = "string"
	TOKEN_INT    = "int"
	TOKEN_BOOL   = "bool"
	TOKEN_BYTE   = "byte"

	NUMBER     = "NUMBER"
	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"

	EOF = "EOF"
	ERR = "ERR"
)

type Token struct {
	Type     TokenType
	Value    string
	Line     int
	Col      int
	Filename string
}

var keywords = map[string]TokenType{
	"string": TOKEN_STRING,
	"int":    TOKEN_INT,
	"bool":   TOKEN_BOOL,
	"byte":   TOKEN_BYTE,
	"float":  NUMBER,
	"fun":    TOKEN_FUN,
	"nil":    TOKEN_NIL,
	"if":     TOKEN_IF,
	"else":   TOKEN_ELSE,
	"for":    TOKEN_FOR,
	"while":  TOKEN_WHILE,
	"return": TOKEN_RETURN,
	"and":    TOKEN_AND,
	"or":     TOKEN_OR,
	"true":   TOKEN_TRUE,
	"false":  TOKEN_FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

func (t Token) String() string {
	return fmt.Sprintf("|Type: %s Value: '%s' Position: %d:%d|", t.Type, t.Value, t.Line, t.Col)
}
