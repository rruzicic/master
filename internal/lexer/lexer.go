package lexer

import (
	"bytes"
	"fmt"
	"interpreter/internal/token"
)

type Lexer struct {
	input    string
	line     int
	col      int
	position int
	ch       byte
	HasError bool
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.line = 1
	return l
}

func (l *Lexer) Tokenize() []token.Token {
	tokens := []token.Token{}
	for {
		l.advance()
		if l.isAtEnd() {
			tokens = append(tokens, l.generateToken(token.EOF))
			break
		}
		l.eatWhitespace()
		switch l.ch {
		case '+':
			tokens = append(tokens, l.generateToken(token.TOKEN_PLUS))
		case '-':
			tokens = append(tokens, l.generateToken(token.TOKEN_MINUS))
		case '*':
			tokens = append(tokens, l.generateToken(token.TOKEN_MUL))
		case '[':
			tokens = append(tokens, l.generateToken(token.TOKEN_LBRACKET))
		case ']':
			tokens = append(tokens, l.generateToken(token.TOKEN_RBRACKET))
		case '(':
			tokens = append(tokens, l.generateToken(token.TOKEN_LPAREN))
		case ')':
			tokens = append(tokens, l.generateToken(token.TOKEN_RPAREN))
		case '{':
			tokens = append(tokens, l.generateToken(token.TOKEN_LCURLY))
		case '}':
			tokens = append(tokens, l.generateToken(token.TOKEN_RCURLY))
		case ';':
			tokens = append(tokens, l.generateToken(token.TOKEN_SEMICOLON))
		case ',':
			tokens = append(tokens, l.generateToken(token.TOKEN_COMMA))
		case '>':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.TOKEN_GTE))
			} else {
				tokens = append(tokens, l.generateToken(token.TOKEN_GT))
			}
		case '<':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.TOKEN_LTE))
			} else {
				tokens = append(tokens, l.generateToken(token.TOKEN_LT))
			}
		case '=':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.TOKEN_EQUAL))
			} else {
				tokens = append(tokens, l.generateToken(token.TOKEN_ASSIGN))
			}
		case '!':
			if l.match('=') {
				tokens = append(tokens, l.generateToken(token.TOKEN_NOT_EQUAL))
			} else {
				tokens = append(tokens, l.generateToken(token.TOKEN_BANG))
			}
		case '/':
			if l.match('/') {
				l.comment()
				tokens = append(tokens, l.generateToken(token.COMMENT))
			} else {
				tokens = append(tokens, l.generateToken(token.TOKEN_DIV))
			}
		case '"':
			str, err := l.sstring()
			if err != nil {
				tokens = append(tokens, l.generateTokenWithValue(token.ERR, err.Error()))
			} else {
				tokens = append(tokens, l.generateTokenWithValue(token.STRING, str))
			}
		case 0:
			tokens = append(tokens, l.generateToken(token.EOF))
		default:
			if isDigit(l.ch) {
				tokens = append(tokens, l.number())
			} else if isAlphaNum(l.ch) {
				tokens = append(tokens, l.identifier())
			} else {
				tokens = append(tokens, l.generateToken(token.ERR))
				l.HasError = true
			}
		}
		if l.isAtEnd() {
			break
		}
	}
	return tokens
}

func (l *Lexer) sstring() (string, error) {
	var buffer bytes.Buffer
	for l.peek() != '"' {
		if l.isAtEnd() {
			return "", fmt.Errorf("%d:%d: unterminated string", l.line, l.col)
		}
		l.advance()
		buffer.WriteByte(l.ch)
	}
	l.advance()
	return buffer.String(), nil
}

func (l *Lexer) number() token.Token {
	var buffer bytes.Buffer
	buffer.WriteByte(l.ch)
	hadDot := false
	for isDigit(l.peek()) || l.peek() == '.' {
		if l.peek() == '.' && hadDot {
			return l.generateToken(token.ERR)
		}
		if l.peek() == '.' {
			hadDot = true
		}
		l.advance()
		buffer.WriteByte(l.ch)
	}
	return l.generateTokenWithValue(token.NUMBER, buffer.String())
}

func (l *Lexer) identifier() token.Token {
	var buffer bytes.Buffer
	buffer.WriteByte(l.ch)
	for isAlphaNum(l.peek()) {
		l.advance()
		buffer.WriteByte(l.ch)
	}
	tok := token.LookupIdent(buffer.String())
	if tok == token.IDENTIFIER {
		return l.generateTokenWithValue(token.IDENTIFIER, buffer.String())
	} else if tok == token.TOKEN_TRUE {
		return l.generateTokenWithValue(token.TOKEN_TRUE, "true")
	} else if tok == token.TOKEN_FALSE {
		return l.generateTokenWithValue(token.TOKEN_FALSE, "false")
	} else {
		return l.generateToken(tok)
	}
}

func (l *Lexer) comment() string {
	var buffer bytes.Buffer
	for l.peek() != '\n' {
		l.advance()
		buffer.WriteByte(l.ch)
	}
	return buffer.String()
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
	if l.position >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}
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
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.advance()
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.position >= len(l.input)
}
