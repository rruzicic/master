package parser

import (
	"interpreter/lexer"
	"testing"
)

func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{`int a=5;
			int b = 3;
		`, "x", 5},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		// t.Log(l.Tokenize())
		program := p.ParseProgram()
		t.Log(program)
		t.Log(p.errors)
		// t.Log(fmt.Sprintf("%T", program.Statements[0]))

		// checkParserErrors(t, p)
		// if len(program.Statements) != 1 {
		// 	t.Fatalf("program.Statements does not contain 1 statements. got=%d",
		// 		len(program.Statements))
		// }
		// stmt := program.Statements[0]
		// if !testLetStatement(t, stmt, tt.expectedIdentifier) {
		// 	return
		// }
		// val := stmt.(*ast.LetStatement).Value
		// if !testLiteralExpression(t, val, tt.expectedValue) {
		// 	return
		// }
	}
}
