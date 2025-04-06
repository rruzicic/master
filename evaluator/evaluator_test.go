package evaluator_test

import (
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
)

func TestInfixExpressionEvaluation(t *testing.T) {
	testCases := []struct {
		input       string
		returnType  string
		returnValue string
	}{
		{
			input:       "true == true;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
		{
			input:       "false == false;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
		{
			input:       "true == false;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "false",
		},
		{
			input:       "false == true;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "false",
		},
		{
			input:       "2+3;",
			returnType:  object.INTEGER_OBJ,
			returnValue: "5",
		},
		{
			input:       "2+3 == 5;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
		{
			input:       "10 / (2+3) == 2;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
		{
			input:       "0 == 0;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
	}
	for i, tC := range testCases {
		l := lexer.New(tC.input)
		if l.HasError {
			t.Fatalf("tests[%d]: lexer errors found", i)

		}
		p := parser.New(l)
		if len(p.Errors()) != 0 {
			t.Fatalf("tests[%d]: parse errors found: %s", i, p.Errors())
		}
		prog := p.ParseProgram()
		env := object.NewEnvironment()
		if len(prog.Statements) != 1 {
			t.Fatalf("tests[%d]: expected 1 statement, got %d", i, len(prog.Statements))
		}
		eval := evaluator.Eval(prog.Statements[0], env)
		if eval.Type() != object.ObjectType(tC.returnType) {
			t.Fatalf("tests[%d]: expected %s object, got %s", i, object.ObjectType(tC.returnType), eval.Type())
		}
		if eval.Inspect() != tC.returnValue {
			t.Fatalf("tests[%d]: expected %s value, got %s", i, tC.returnValue, eval.Inspect())
		}

	}
}
