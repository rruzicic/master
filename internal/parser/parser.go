package parser

import (
	"fmt"
	"interpreter/internal/ast"
	"interpreter/internal/lexer"
	"interpreter/internal/token"
	"strconv"
	"strings"
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	LOGICAL
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.TOKEN_EQUAL:     ASSIGN,
	token.TOKEN_NOT_EQUAL: ASSIGN,
	token.TOKEN_AND:       LOGICAL,
	token.TOKEN_OR:        LOGICAL,
	token.TOKEN_LT:        LESSGREATER,
	token.TOKEN_GT:        LESSGREATER,
	token.TOKEN_LTE:       LESSGREATER,
	token.TOKEN_GTE:       LESSGREATER,
	token.TOKEN_PLUS:      SUM,
	token.TOKEN_MINUS:     SUM,
	token.TOKEN_MUL:       PRODUCT,
	token.TOKEN_DIV:       PRODUCT,
	token.TOKEN_LPAREN:    CALL,
	token.TOKEN_LBRACKET:  INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	tokens    []token.Token
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	// TODO: use channels or directly call NextToken, don't use arrays
	p.curToken = p.peekToken
	if p.curToken.Type == token.ERR {
		p.errors = append(p.errors, "received an error token")
	}
	if len(p.tokens) == 1 {
		p.peekToken = p.tokens[0]
	} else {
		p.peekToken, p.tokens = p.tokens[0], p.tokens[1:]
	}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curIsTypeToken() bool {
	return p.curToken.Type == token.TOKEN_INT || p.curToken.Type == token.TOKEN_STRING || p.curToken.Type == token.TOKEN_BOOL || p.curToken.Type == token.TOKEN_BYTE || p.curToken.Type == token.TOKEN_FLOAT
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.tokens = l.Tokenize()
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberExpression)
	p.registerPrefix(token.STRING, p.parseStringExpression)
	p.registerPrefix(token.TOKEN_TRUE, p.parseBoolExpression)
	p.registerPrefix(token.TOKEN_FALSE, p.parseBoolExpression)
	p.registerPrefix(token.TOKEN_BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TOKEN_MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TOKEN_LBRACKET, p.parseArrayExpression)
	p.registerPrefix(token.TOKEN_LPAREN, p.parseGroupExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.TOKEN_PLUS, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_MINUS, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_MUL, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_DIV, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_GT, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_LT, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_GTE, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_LTE, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_OR, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_AND, p.parseInfixExpression)
	p.registerInfix(token.TOKEN_LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.TOKEN_LPAREN, p.parseCallExpression)
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.TOKEN_VAR:
		return p.parseVarStatement()
	case token.IDENTIFIER:
		if p.peekToken.Type == token.TOKEN_ASSIGN {
			return p.parseVarStatement()
		}
	case token.TOKEN_WHILE:
		return p.parseWhileStatement()
	case token.TOKEN_IF:
		return p.parseIfStatement()
	case token.TOKEN_RETURN:
		return p.parseReturnStatement()
	case token.TOKEN_LCURLY:
		return p.parseBlockStatement()
	case token.TOKEN_FUN:
		return p.parseFunctionDefinition()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}
	if p.curToken.Type == token.TOKEN_VAR {
		p.nextToken()
	}
	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors, "expected identifier")
		return nil
	}
	stmt.Identifier = &ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
	p.nextToken()
	if p.curToken.Type == token.TOKEN_ASSIGN {
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
		p.nextToken()
		return stmt
	} else if p.curToken.Type == token.TOKEN_SEMICOLON {
		return stmt
	} else {
		p.errors = append(p.errors, "expected = or ; got neither")
		return nil
	}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	stmt.Body = *p.parseBlockStatement()
	return stmt
}

// TODO: implement IF ... ELSE IF .... ELSE
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	stmt.Body = p.parseBlockStatement()
	if p.peekToken.Type == token.TOKEN_ELSE {
		p.nextToken()
		stmt.Alternative = p.parseBlockStatement()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.TOKEN_SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseFunctionDefinition() ast.Statement {
	stmt := &ast.FunctionStatement{
		Token: p.curToken,
	}
	p.nextToken()
	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors, "function definition missing identifier")
		return nil
	}
	stmt.Identifier = p.curToken
	p.nextToken()

	stmt.ParameterList = p.parseFunctionParameterList()

	p.nextToken()

	if p.curToken.Type != token.TOKEN_LCURLY {
		p.errors = append(p.errors, fmt.Sprintf("expected {, got %s", p.curToken.Type))
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameterList() []ast.IdentifierExpression {
	paramList := []ast.IdentifierExpression{}
	if p.peekToken.Type == token.TOKEN_RPAREN {
		p.nextToken()
		return paramList
	}

	p.nextToken()
	ident := ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
	paramList = append(paramList, ident)
	for p.peekToken.Type == token.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		ident := ast.IdentifierExpression{
			Token: p.curToken,
			Value: p.curToken.Value,
		}
		paramList = append(paramList, ident)
	}

	if !p.expectPeek(token.TOKEN_RPAREN) {
		return nil
	}

	return paramList
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmt := &ast.BlockStatement{Token: p.curToken}
	stmt.Statements = []ast.Statement{}
	p.nextToken()
	for p.curToken.Type != token.TOKEN_RCURLY && p.curToken.Type != token.EOF {
		parsedStmt := p.parseStatement()
		if parsedStmt != nil {
			stmt.Statements = append(stmt.Statements, parsedStmt)
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.TOKEN_SEMICOLON) {
		p.nextToken()
	}
	// TODO: throw error here
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.TOKEN_SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: string(p.curToken.Type), // TODO: set value of token to .Value not just .Type
	}
	prec := LOWEST
	if p, ok := precedences[p.curToken.Type]; ok {
		prec = p
	}
	p.nextToken()
	exp.Right = p.parseExpression(prec)
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: string(p.curToken.Type),
	}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseNumberExpression() ast.Expression {
	if floatExp := p.parseFloatNumber(); floatExp != nil {
		return floatExp
	}
	// parse int
	valueInt, err := strconv.ParseInt(p.curToken.Value, 0, 64)
	if err == nil {
		lit := &ast.IntegerLiteral{Token: p.curToken}
		lit.Value = valueInt
		return lit
	}
	// TODO: parse byte

	msg := fmt.Sprintf("could not parse %q as integer, float or byte, %s", p.curToken.Value, err)
	p.errors = append(p.errors, msg)
	return nil
}

func (p *Parser) parseFloatNumber() ast.Expression {
	if !strings.Contains(p.curToken.Value, ".") {
		return nil
	}
	if valueFloat, err := strconv.ParseFloat(p.curToken.Value, 64); err == nil {
		lit := &ast.FloatLiteral{Token: p.curToken}
		lit.Value = valueFloat
		return lit
	}
	return nil
}

func (p *Parser) parseBoolExpression() ast.Expression {
	lit := &ast.BoolLiteral{Token: p.curToken}
	if p.curToken.Value == "true" {
		lit.Value = true
	} else {
		lit.Value = false
	}
	return lit
}

func (p *Parser) parseStringExpression() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken}
	lit.Value = p.curToken.Value
	return lit
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.TOKEN_RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayExpression() ast.Expression {
	exp := &ast.ArrayLiteral{
		Token: p.curToken,
	}
	if p.peekToken.Type == token.TOKEN_RBRACKET {
		p.nextToken()
		return exp
	}
	p.nextToken()
	exp.Values = append(exp.Values, p.parseExpression(LOWEST))
	for p.peekToken.Type == token.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		exp.Values = append(exp.Values, p.parseExpression(LOWEST))
	}
	p.expectPeek(token.TOKEN_RBRACKET)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	p.expectPeek(token.TOKEN_RBRACKET)
	return exp
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	funIdent, ok := left.(*ast.IdentifierExpression)
	if !ok {
		p.errors = append(p.errors, "invalid function identifier")
		return nil
	}
	exp := &ast.CallExpression{
		Token:             p.curToken,
		FunctionIdentifer: funIdent,
	}
	if p.peekToken.Type == token.TOKEN_RPAREN {
		p.nextToken()
		return exp
	}
	p.nextToken()
	exp.Parameters = append(exp.Parameters, p.parseExpression(LOWEST))
	for p.peekToken.Type == token.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		exp.Parameters = append(exp.Parameters, p.parseExpression(LOWEST))
	}
	if p.peekToken.Type != token.TOKEN_RPAREN {
		p.errors = append(p.errors, "expected )")
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) Errors() []string {
	return p.errors
}
