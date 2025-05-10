package object_test

import (
	"interpreter/internal/object"
	"testing"
)

func TestArrayInspect(t *testing.T) {
	testCases := []struct {
		array          object.Array
		expectedOutput string
	}{
		{
			array:          object.Array{Elements: []object.Object{&object.Boolean{Value: true}, &object.String{Value: "abc"}, &object.Integer{Value: 42}}},
			expectedOutput: "[true, abc, 42]",
		},
	}
	for i, tC := range testCases {
		obj := &tC.array
		testObjectInspect(t, i, obj, tC.expectedOutput)
	}
}

func TestErrorInspect(t *testing.T) {
	testCases := []struct {
		err            object.Error
		expectedOutput string
	}{
		{
			err: object.Error{
				Error: "this is an error",
			},
			expectedOutput: "ERROR: this is an error",
		},
	}
	for i, tC := range testCases {
		obj := &tC.err
		testObjectInspect(t, i, obj, tC.expectedOutput)
	}
}

func TestFloatInspect(t *testing.T) {
	testObjectInspect(t, 0, &object.Float{Value: 5.1}, "5.100000")
}

func TestStdfuncInspect(t *testing.T) {
	testObjectInspect(t, 0, &object.StdFunction{}, "<std fun>")
}
func TestFunInspect(t *testing.T) {
	testObjectInspect(t, 0, &object.Function{}, "<fun>")
}

func TestReturnValue(t *testing.T) {
	testObjectInspect(t, 0, &object.ReturnValue{Value: &object.Boolean{true}}, "true")
}

func TestBoolInspect(t *testing.T) {
	testObjectInspect(t, 0, &object.Boolean{Value: true}, "true")
}
func TestNullInspect(t *testing.T) {
	testObjectInspect(t, 0, &object.Null{}, "null")
}

func testObjectInspect(t *testing.T, tstNum int, obj object.Object, expected string) {
	result := obj.Inspect()
	if result != expected {
		t.Fatalf("test[%d]: expected %s , got %s", tstNum, expected, result)

	}
}
