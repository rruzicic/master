package lexer

import (
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
		case '/': // this could also be a comment
			tokens = append(tokens, l.generateToken(token.DIV))
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
			if l.peek('=') {
				l.advance()
				tokens = append(tokens, l.generateToken(token.GTE))
			} else {
				tokens = append(tokens, l.generateToken(token.GT))
			}
		case '<':
			if l.peek('=') {
				l.advance()
				tokens = append(tokens, l.generateToken(token.LTE))
			} else {
				tokens = append(tokens, l.generateToken(token.LT))
			}
		case '=':
			if l.peek('=') {
				l.advance()
				tokens = append(tokens, l.generateToken(token.EQUAL))
			} else {
				tokens = append(tokens, l.generateToken(token.ASSIGN))
			}
		case '!':
			if l.peek('=') {
				l.advance()
				tokens = append(tokens, l.generateToken(token.NOT_EQUAL))
			} else {
				tokens = append(tokens, l.generateToken(token.BANG))
			}

		default:
			l.eatWhitespace()
		}
	}
	tokens = append(tokens, l.generateToken(token.EOF))
	return tokens
}

func (l *Lexer) generateTokenWithValue(typez string, value string) token.Token {
	return token.Token{
		Type:  typez,
		Value: value,
		Line:  l.line,
		Col:   l.col,
	}

}

func (l *Lexer) generateToken(typez string) token.Token {
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

func (l *Lexer) peek(ch byte) bool {
	if l.isAtEnd() {
		return false
	}
	return l.input[l.position] == ch
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
