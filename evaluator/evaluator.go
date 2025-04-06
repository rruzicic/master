package evaluator

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)
	// literals
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.ArrayLiteral:
	case *ast.BoolLiteral:
		fmt.Println(node)
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// identifier
	case *ast.IdentifierExpression:

	// expressions
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left
		}
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalInfixExpression(left, right, node.Operator)
	case *ast.PrefixExpression:
	case *ast.CallExpression:
	case *ast.IndexExpression:

		// statements
	case *ast.BlockStatement:
	case *ast.FunctionStatement:
	case *ast.IfStatement:
	case *ast.ReturnStatement:
	case *ast.VarStatement:
	case *ast.WhileStatement:
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	}
	return nil
}

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var ret object.Object
	for _, s := range node.Statements {
		ret = Eval(s, env)
		switch res := ret.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return ret
}

func evalInfixExpression(left object.Object, right object.Object, operator string) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		l := left.(*object.Integer).Value
		r := right.(*object.Integer).Value
		switch operator {
		case "+":
			return &object.Integer{Value: l + r}
		case "-":
			return &object.Integer{Value: l - r}
		case "/":
			return &object.Integer{Value: l / r}
		case "*":
			return &object.Integer{Value: l * r}
		case ">":
			return nativeBoolToBooleanObject(l > r)
		case "<":
			return nativeBoolToBooleanObject(l < r)
		case "==":
			return nativeBoolToBooleanObject(l == r)
		case "!=":
			return nativeBoolToBooleanObject(l != r)
		default:
			return &object.Error{Error: "unknown operator: " + operator}
		}
	case (left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ) ||
		(right.Type() == object.FLOAT_OBJ && left.Type() == object.INTEGER_OBJ) ||
		(left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ):
		var l float64
		var r float64
		if left.Type() == object.INTEGER_OBJ {
			l = float64(left.(*object.Integer).Value)
		} else {
			l = left.(*object.Float).Value
		}

		if right.Type() == object.INTEGER_OBJ {
			r = float64(right.(*object.Integer).Value)
		} else {
			r = right.(*object.Float).Value
		}
		switch operator {

		case "+":
			return &object.Float{Value: l + r}
		case "-":
			return &object.Float{Value: l - r}
		case "/":
			return &object.Float{Value: l / r}
		case "*":
			return &object.Float{Value: l * r}
		case ">":
			return nativeBoolToBooleanObject(l > r)
		case "<":
			return nativeBoolToBooleanObject(l < r)
		case "==":
			return nativeBoolToBooleanObject(l == r)
		case "!=":
			return nativeBoolToBooleanObject(l != r)
		default:
			return &object.Error{Error: "unknown operator: " + operator}
		}
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return &object.Error{Error: "type mismatch"}
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		if operator != "+" {
			return &object.Error{Error: "operator other than + not allowed for strings"}
		}
		l := left.(*object.String).Value
		r := right.(*object.String).Value
		return &object.String{Value: l + r}
	default:
		return &object.Error{Error: "unknown error "}
	}
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
