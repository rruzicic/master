package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"strings"
	"testing"
)

func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var a = 5;", "a", 5},
		{"var b = a;", "b", "a"},
		{"var x = false;", "x", false},
		{"var y = 571.1;", "y", 571.1},
		{"var z = [ 1 , 2 , 3 ];", "z", []int{1, 2, 3}},
		{"var h = [ 1 ];", "h", []int{1}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
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

// TODO: more test cases(spaces, tabs, weird charcters etc)
func TestStringLiteralExpression(t *testing.T) {
	input := `var pera = "mira";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.VarStatement)
	literal, ok := stmt.Value.(*ast.StringLiteral)
	if stmt.Identifier.Value != "pera" {
		t.Fatalf("var stmt identifier not pera. got=%s", stmt.Identifier.Value)
	}
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Value)
	}

	if literal.Value != "mira" {
		t.Errorf("literal.Value not %q. got=%q", "mira", literal.Value)
	}

}

// TODO: more test cases(and different return types: string, bool, full expressions)
func TestReturnStatement(t *testing.T) {
	input := `return 2+2;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ReturnStatement")
	}
	if stmt.Value.String() != "(2 + 2)" {
		t.Fatalf("stmt value not '(2+2)'. got=%s", stmt.Value)
	}

}

// TODO: test only if
func TestIfElseStatement(t *testing.T) {
	input := `
	if (pera == 3) {
		var jova = 5;
	} else {
		var djoka = 8;
	}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("exp not *ast.IfStatement. got=%T", stmt)
	}

	if stmt.Condition.String() != "(pera == 3)" {
		t.Fatalf("condition not (pera == 3). got=%s", stmt.Condition)
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("number of body statements not 1. got=%d", len(stmt.Body.Statements))
	}

	if stmt.Body.Statements[0].String() != "var jova = 5;" {
		t.Fatalf("body statement not 'var jova = 5;'. got='%s'", stmt.Body.Statements[0])
	}

	if len(stmt.Alternative.Statements) != 1 {
		t.Fatalf("number of alternative statements not 1. got=%d", len(stmt.Alternative.Statements))
	}

	if strings.Trim(stmt.Alternative.Statements[0].String(), " \n\t") != "var djoka = 8;" {
		t.Fatalf("alternative statement not 'var djoka = 8;'. got='%s'", strings.Trim(stmt.Alternative.Statements[0].String(), " \n\t"))
	}
}

// TODO: expand block statement tests when new stmts/expr get implemented
// TODO: nested block statements
func TestBlockStatement(t *testing.T) {
	tests := []struct {
		input                 string
		expectedNumberOfStmts int
	}{
		// TODO: add dedicated test for IndexExpression
		{`
		{
			var a = pera(1);
			pera(1,2,3);
			var c = djoka(634, 234, 234);
			var d=a;
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

// TODO: more tests(with nesting)
func TestWhileStatement(t *testing.T) {
	input := `
	while (zika > 3) {
		var a = a + 1;
		a + 1;
		b - 5;
	}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("exp not *ast.WhileStatement. got=%T", stmt)
	}

}

// TODO: check return type, parameters, function name
func TestFunctionDefinition(t *testing.T) {

	inputs := []string{`
	fun x(int a, int b) int {
		var b = b + 1;
		a = a + 1;
		return a;
	}
	`,
		`
	fun a() int { return 1; } a();
	`,
	}
	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt, ok := program.Statements[0].(*ast.FunctionStatement)
		if !ok {
			t.Fatalf("exp not *ast.FunctionStatement. got=%T", stmt)
		}
	}
}

// TODO: consider merging prefix and infix expression tests
func TestInfixExpressions(t *testing.T) { // TODO: add more tests
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"5 + 10 / 2;", "(5 + (10 / 2));"},
		{"pera();", "(pera());"},
		{"pera(1,2,3);", "(pera(1, 2, 3));"},
		{"a(1);", "(a(1));"},
		{"b = 5 + 2 / 3 - 6 * 9 - a(1);", "b = (((5 + (2 / 3)) - (6 * 9)) - (a(1)));"},
		{"var j = 9123 - a[81] * (12 - 3);", "var j = (9123 - ((a[81]) * (12 - 3)));"},
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

// TODO: add more tests after adding support for group expressions
func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"-5;", "(-5);"},
		{"-5 + 2;", "((-5) + 2);"},
		{"-(5 + 2);", "(-(5 + 2));"},
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
	if varStmt.Identifier == nil {
		t.Errorf("varStmt.Name is nil")
		return false

	}
	if varStmt.Identifier.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Identifier.Value)
		return false
	}
	if varStmt.Identifier.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Identifier)
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanLiteral(t, exp, v)
	case []int:
		return testIntegerArrayLiteral(t, exp, v)
	// case []int64:
	// 	return testIntegerArrayLiteral(t, exp, []int(v))
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float32:
		return testFloatLiteral(t, exp, float64(v))
	case float64:
		return testFloatLiteral(t, exp, float64(v))
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

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	float, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}
	if float.Value != value {
		t.Errorf("fl.Value not %g. got=%g", value, float.Value)
		return false
	}
	if float.TokenLiteral() != fmt.Sprintf("%g", value) {
		t.Errorf("fl.TokenLiteral not %g. got=%s", value, float.TokenLiteral())
		return false
	}
	return true
}

func testIntegerArrayLiteral(t *testing.T, al ast.Expression, value []int) bool {
	arr, ok := al.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("al not *ast.ArrayLiteral. got=%T", al)
		return false
	}
	for i, v := range arr.Values {
		il, ok := v.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("il[%d] not *ast.IntegerLiteral. got=%T", i, il)
			return false
		}
	}
	if arr.String() != fmt.Sprintf("%d", value) {
		t.Errorf("al.TokenLiteral not %d. got=%s", value, arr.String())
		return false
	}
	return true
}
