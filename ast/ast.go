package ast

import (
	"bytes"
	"interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var ret bytes.Buffer
	for _, s := range p.Statements {
		ret.WriteString(s.String())
	}
	return ret.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *InfixExpression) String() string {
	var ret bytes.Buffer
	ret.WriteString("(")
	ret.WriteString(ie.Left.String())
	ret.WriteString(" " + ie.Operator + " ")
	ret.WriteString(ie.Right.String())
	ret.WriteString(")")
	return ret.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Value }
func (pe *PrefixExpression) String() string {
	var ret bytes.Buffer
	ret.WriteString("(")
	ret.WriteString(pe.Operator)
	ret.WriteString(pe.Right.String())
	ret.WriteString(")")
	return ret.String()
}

type IdentifierExpression struct {
	Token token.Token
	Value string
}

func (ie *IdentifierExpression) expressionNode()      {}
func (ie *IdentifierExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IdentifierExpression) String() string       { return ie.Value }

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Value }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }
