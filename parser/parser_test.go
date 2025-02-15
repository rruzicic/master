package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"int a = 5;", "a", 5},
		{"int b = a;", "b", "a"},
		{"bool x = false;", "x", false},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `string pera = "mira";`
	l := lexer.New(input)
	// t.Log(l.Tokenize())
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.VarStatement)
	literal, ok := stmt.Value.(*ast.StringLiteral)
	if stmt.Name.Value != "pera" {
		t.Fatalf("var stmt identifier not pera. got=%s", stmt.Name.Value)
	}
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Value)
	}

	if literal.Value != "mira" {
		t.Errorf("literal.Value not %q. got=%q", "mira", literal.Value)
	}

}

// TODO: expand block statement tests when new stmts/expr get implemented
func TestBlockStatement(t *testing.T) {
	tests := []struct {
		input                 string
		expectedNumberOfStmts int
	}{
		{`
		{
			int a = 5;
			bool b = true;
		}
		`, 1},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != tt.expectedNumberOfStmts {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", tt.expectedNumberOfStmts, len(program.Statements))
		}
		bs, ok := program.Statements[0].(*ast.BlockStatement)
		if !ok {
			t.Errorf("exp not *ast.BlockStatement. got=%T", bs)
			return
		}

	}
}

// TODO: consider merging prefix and infix expression tests
func TestInfixExpressions(t *testing.T) { // TODO: add more tests
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"5 + 10 / 2;", "(5 + (10 / 2))"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		if program.Statements[0].String() != tt.expectedOutput {
			t.Errorf("exp not %s. got=%s", tt.expectedOutput, program.Statements[0])
			return
		}

	}
}

func TestPrefixExpressions(t *testing.T) { // TODO: add more tests after adding support for group expressions
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"-5;", "(-5)"},
		{"-5 + 2;", "((-5) + 2)"},
		{"-(5 + 2);", "(-(5 + 2))"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		if program.Statements[0].String() != tt.expectedOutput {
			t.Errorf("exp not %s. got=%s", tt.expectedOutput, program.Statements[0])
			return
		}

	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}
	if varStmt.Name == nil {
		t.Errorf("varStmt.Name is nil")
		return false

	}
	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}
	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanLiteral(t, exp, v)
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bl, ok := exp.(*ast.BoolLiteral)
	if !ok {
		t.Errorf("exp not *ast.BoolLiteral. got=%T", exp)
		return false
	}
	if bl.Value != value {
		t.Errorf("bl.Value not %t. got=%t", value, bl.Value)
		return false
	}
	if bl.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bl.TokenLiteral not %t. got=%s", value, bl.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.IdentifierExpression)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}
