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
		returnType  object.ObjectType
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
		{
			input:       "1.1 + 5;",
			returnType:  object.FLOAT_OBJ,
			returnValue: "6.100000",
		},
		{
			input:       "1.1 + 5.2;",
			returnType:  object.FLOAT_OBJ,
			returnValue: "6.300000",
		},
		{
			input:       "1.1 == 1.1;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
		{
			input:       "true or false;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
		{
			input:       "true and false;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "false",
		},
	}
	for i, tC := range testCases {
		eval := evaluate(t, i, tC.input)
		checkTypeAndValue(t, i, eval, tC.returnType, tC.returnValue)
	}
}

func TestVarEvaluation(t *testing.T) {
	testCases := []struct {
		input       string
		returnType  object.ObjectType
		returnValue string
	}{
		{
			input:       "var a = 5; a + 1;",
			returnType:  object.INTEGER_OBJ,
			returnValue: "6",
		},
		{
			input:       "var a = 5; var b = 1; a + b;",
			returnType:  object.INTEGER_OBJ,
			returnValue: "6",
		},
		{
			input:       "var b = true; true == b;",
			returnType:  object.BOOLEAN_OBJ,
			returnValue: "true",
		},
	}
	for i, tC := range testCases {
		eval := evaluate(t, i, tC.input)
		checkTypeAndValue(t, i, eval, tC.returnType, tC.returnValue)
	}
}

func TestFunctionEvaluation(t *testing.T) {
	testCases := []struct {
		input       string
		returnType  object.ObjectType
		returnValue string
	}{
		{
			input:       "fun a(x, y) { return x+y; } a(2,3);",
			returnType:  object.INTEGER_OBJ,
			returnValue: "5",
		},
		{
			input:       "fun a(b) { return b; } a(\"abc\");",
			returnType:  object.STRING_OBJ,
			returnValue: "abc",
		},
		{
			input:       "fun a() { return 1; } a();",
			returnType:  object.INTEGER_OBJ,
			returnValue: "1",
		},
	}
	for i, tC := range testCases {
		eval := evaluate(t, i, tC.input)
		checkTypeAndValue(t, i, eval, tC.returnType, tC.returnValue)
	}
}

func TestStdFunctions(t *testing.T) {
	testCases := []struct {
		input       string
		returnType  object.ObjectType
		returnValue string
	}{
		{
			input:       `len("hello");`,
			returnType:  object.INTEGER_OBJ,
			returnValue: "5",
		},
		{
			input:       `len(["hello", 1, 5.2]);`,
			returnType:  object.INTEGER_OBJ,
			returnValue: "3",
		},
	}
	for i, tC := range testCases {
		eval := evaluate(t, i, tC.input)
		checkTypeAndValue(t, i, eval, tC.returnType, tC.returnValue)
	}
}

func evaluate(t *testing.T, testNum int, input string) object.Object {
	l := lexer.New(input)
	if l.HasError {
		t.Fatalf("tests[%d]: lexer errors found", testNum)
	}
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("tests[%d]: parse errors found: %s", testNum, p.Errors())
	}
	env := object.NewEnvironment()
	return evaluator.Eval(prog, env)
}

func checkTypeAndValue(t *testing.T, testNum int, eval object.Object, returnType object.ObjectType, returnValue string) {
	if eval.Type() != returnType {
		t.Fatalf("tests[%d]: expected %s object, got %s", testNum, object.ObjectType(returnType), eval.Type())
	}
	if eval.Inspect() != returnValue {
		t.Fatalf("tests[%d]: expected %s value, got %s", testNum, returnValue, eval.Inspect())
	}
}
