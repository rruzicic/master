package lexer

import (
	"bytes"
	"errors"
	"interpreter/token"
)

type Lexer struct {
	input    string
	line     int
	col      int
	position int
	ch       byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.line = 1
	return l
}

func (l *Lexer) Tokenize() []token.Token {
	tokens := []token.Token{}
	for {
		if l.isAtEnd() {
			break
		}
		l.advance()
		switch l.ch {
		case '+':
			tokens = append(tokens, l.generateToken(token.PLUS))
		case '-':
			tokens = append(tokens, l.generateToken(token.MINUS))
		case '*':
			tokens = append(tokens, l.generateToken(token.MUL))
		case '[':
			tokens = append(tokens, l.generateToken(token.LBRACKET))
		case ']':
			tokens = append(tokens, l.generateToken(token.RBRACKET))
		case '(':
			tokens = append(tokens, l.generateToken(token.LPAREN))
		case ')':
			tokens = append(tokens, l.generateToken(token.RPAREN))
		case '{':
			tokens = append(tokens, l.generateToken(token.LCURLY))
		case '}':
			tokens = append(tokens, l.generateToken(token.RCURLY))
		case ';':
			tokens = append(tokens, l.generateToken(token.SEMICOLON))
		case '>':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.GTE))
			} else {
				tokens = append(tokens, l.generateToken(token.GT))
			}
		case '<':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.LTE))
			} else {
				tokens = append(tokens, l.generateToken(token.LT))
			}
		case '=':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.EQUAL))
			} else {
				tokens = append(tokens, l.generateToken(token.ASSIGN))
			}
		case '!':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.NOT_EQUAL))
			} else {
				tokens = append(tokens, l.generateToken(token.BANG))
			}
		case '/': // this could also be a comment
			tokens = append(tokens, l.generateToken(token.DIV))
		case '"':
			str, err := l.sstring()
			if err != nil {
				tokens = append(tokens, l.generateTokenWithValue(token.ERR, err.Error()))
			} else {
				tokens = append(tokens, l.generateTokenWithValue(token.STRING, str))
			}
		default:
			if isDigit(l.ch) {
				tokens = append(tokens, l.number())
			} else if isAlphaNum(l.ch) {
				tokens = append(tokens, l.identifier())
			} else {
				l.eatWhitespace()
			}
		}
	}
	tokens = append(tokens, l.generateToken(token.EOF))
	return tokens
}

func (l *Lexer) sstring() (string, error) {
	var buffer bytes.Buffer
	for l.peek() != '"' {
		if l.isAtEnd() {
			return "", errors.New("unterminated string")
		}
		l.advance()
		buffer.WriteByte(l.ch)
	}
	l.advance()
	return buffer.String(), nil
}

// TODO: implement float and rename func
func (l *Lexer) number() token.Token {
	var buffer bytes.Buffer
	for isDigit(l.ch) {
		buffer.WriteByte(l.ch)
		l.advance()
	}
	return l.generateTokenWithValue(token.FLOAT, buffer.String())
}

func (l *Lexer) identifier() token.Token {
	var buffer bytes.Buffer
	buffer.WriteByte(l.ch)
	for isAlphaNum(l.peek()) {
		l.advance()
		buffer.WriteByte(l.ch)
	}
	// l.advance()
	tok := token.LookupIdent(buffer.String())
	if tok == token.IDENTIFIER {
		return l.generateTokenWithValue(token.IDENTIFIER, buffer.String())
	} else {
		return l.generateToken(tok)
	}
}

func isAlphaNum(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) generateTokenWithValue(typez token.TokenType, value string) token.Token {
	return token.Token{
		Type:  typez,
		Value: value,
		Line:  l.line,
		Col:   l.col,
	}
}

func (l *Lexer) generateToken(typez token.TokenType) token.Token {
	return token.Token{
		Type:  typez,
		Value: "",
		Line:  l.line,
		Col:   l.col,
	}
}

func (l *Lexer) advance() {
	l.ch = l.input[l.position]
	l.position += 1
	l.col += 1
	if l.ch == '\n' {
		l.line += 1
		l.col = 0
	}
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.input[l.position]
}

func (l *Lexer) match(ch byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.input[l.position] == ch {
		l.advance()
		return true
	}
	return false
}

func (l *Lexer) eatWhitespace() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			break
		}
		l.advance()
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.position >= len(l.input)
}
