package token

import "fmt"

type TokenType string

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
	"string": STRINGT,
	"int":    INT,
	"bool":   BOOL,
	"byte":   BYTE,
	"float":  FLOAT,
	"fun":    FUN,
	"nil":    NIL,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"while":  WHILE,
	"return": RETURN,
	"and":    AND,
	"or":     OR,
	"true":   TRUE,
	"false":  FALSE,
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
