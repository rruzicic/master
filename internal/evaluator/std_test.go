package evaluator_test

import (
	"interpreter/internal/object"
	"testing"
)

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
		{
			input:       `len(5.2);`,
			returnType:  object.INTEGER_OBJ,
			returnValue: "0",
		},
	}
	for i, tC := range testCases {
		eval := evaluate(t, i, tC.input)
		checkTypeAndValue(t, i, eval, tC.returnType, tC.returnValue)
	}
}
