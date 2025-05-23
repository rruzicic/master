package ast

import (
	"bytes"
	"interpreter/internal/token"
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
	Type  token.Token
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
		return es.Expression.String() + ";"
	}
	return ""
}

type FunctionStatement struct {
	Token         token.Token // token.TOKEN_FUNCTION
	Identifier    token.Token
	ParameterList []IdentifierExpression
	Body          *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) TokenLiteral() string {
	return fs.Identifier.Value
}
func (fs *FunctionStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString("func ")
	buf.WriteString(fs.Identifier.Value)
	buf.WriteString("(")
	for i, p := range fs.ParameterList {
		buf.WriteString(p.String())
		if len(fs.ParameterList) > i+1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(") ")
	buf.WriteString(fs.Body.String())
	return buf.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FloatLiteral) String() string       { return fl.Token.Value }

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BoolLiteral) expressionNode() {}
func (bl *BoolLiteral) TokenLiteral() string {
	// TODO: set value of Token.Value while scanning tokens and remove code below
	if bl.Value {
		return "true"
	} else {
		return "false"
	}
}
func (bl *BoolLiteral) String() string {
	// TODO: set value of Token.Value while scanning tokens and remove code below
	if bl.Value {
		return "true"
	} else {
		return "false"
	}
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) String() string       { return sl.Token.Value }

// TODO: add byte and float literal nodes
type VarStatement struct {
	Token      token.Token // type token
	Identifier *IdentifierExpression
	Value      Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Value }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	if vs.Token.Type != token.IDENTIFIER {
		out.WriteString(string(vs.Token.Type) + " ")
	}
	out.WriteString(vs.Identifier.String())
	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // token.TOKEN_LCURLY token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString("\n\t" + s.String())
	}
	out.WriteString("\n")
	return out.String()
}

type WhileStatement struct {
	Token     token.Token // token.TOKEN_WHILE token
	Condition Expression
	Body      BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Value }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ws.Token.Value)
	out.WriteString("while (")
	out.WriteString(ws.Condition.String())
	out.WriteString(")")
	out.WriteString(" {")
	out.WriteString(ws.Body.String())
	out.WriteString("}")
	return out.String()
}

type IfStatement struct {
	Token       token.Token // token.TOKEN_WHILE token
	Condition   Expression
	Body        *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Value }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.Token.Value)
	out.WriteString("(")
	out.WriteString(is.Condition.String())
	out.WriteString(")")
	out.WriteString(is.Body.String())
	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}
	return out.String()
}

type ReturnStatement struct {
	Token token.Token // token.TOKEN_RETURN token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Value }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Token.Value)
	out.WriteString(rs.Value.String())
	return out.String()
}

type ArrayLiteral struct {
	Token  token.Token
	Values []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Value }
func (al *ArrayLiteral) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, e := range al.Values {
		buf.WriteString(e.String())
		if len(al.Values) > i+1 {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IndexExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(ie.Left.String())
	buf.WriteString("[")
	buf.WriteString(ie.Index.String())
	buf.WriteString("]")
	buf.WriteString(")")
	return buf.String()
}

type CallExpression struct {
	Token             token.Token // ( token
	FunctionIdentifer Expression
	Parameters        []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(ce.FunctionIdentifer.String())
	buf.WriteString("(")
	for i, p := range ce.Parameters {
		buf.WriteString(p.String())
		if len(ce.Parameters) > i+1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")")
	buf.WriteString(")")
	return buf.String()
}
