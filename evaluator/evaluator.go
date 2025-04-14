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
		elements := evalParameters(node.Values, env)
		if len(elements) == 1 && elements[0].Type() == object.ERROR_OBJ {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.BoolLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// identifier
	case *ast.IdentifierExpression:
		if val, ok := env.Get(node.Value); ok {
			return val
		}
		return &object.Error{Error: fmt.Sprintf("identifier not found: " + node.Value)}
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
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		switch node.Operator {
		case "-":
			switch right.Type() {
			case object.FLOAT_OBJ:
				val := right.(*object.Float).Value
				return &object.Float{Value: -val}
			case object.INTEGER_OBJ:
				val := right.(*object.Integer).Value
				return &object.Integer{Value: -val}
			default:
				return &object.Error{Error: fmt.Sprintf("operator - unsuported for %s", right.Type())}
			}
		case "!":
			switch right {
			case TRUE:
				return FALSE
			case FALSE:
				return TRUE
			case NULL:
				return TRUE
			default:
				return FALSE
			}
		}
		return &object.Error{Error: "unsupported prefix operator"}
	case *ast.CallExpression:
		function := Eval(node.FunctionIdentifer, env)
		if function.Type() == object.ERROR_OBJ {
			return function
		}
		params := evalParameters(node.Parameters, env)
		if len(params) == 1 && params[0].Type() == object.ERROR_OBJ {
			return params[0]
		}
		return evalFunction(function, params)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left
		}
		index := Eval(node.Index, env)
		if index.Type() == object.ERROR_OBJ {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.BlockStatement:
		var ret object.Object
		for _, statement := range node.Statements {
			ret = Eval(statement, env)
			if ret != nil {
				if ret.Type() == object.RETURN_VALUE_OBJ || ret.Type() == object.ERROR_OBJ {
					return ret
				}
			}
		}
		return ret
	case *ast.FunctionStatement:
		function := &object.Function{
			Params: node.ParameterList,
			Body:   node.Body,
			Env:    env,
		}
		env.Set(node.Identifier.Value, function)
		return function
	case *ast.IfStatement:
		condition := Eval(node.Condition, env)
		if condition.Type() == object.ERROR_OBJ {
			return condition
		}
		if isTrue(condition) {
			return Eval(node.Body, env)
		} else {
			if node.Alternative != nil {
				return Eval(node.Alternative, env)
			}
		}
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if val.Type() == object.ERROR_OBJ {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.VarStatement:
		value := Eval(node.Value, env)
		if value.Type() == object.ERROR_OBJ {
			return value
		}
		env.Set(node.Identifier.Value, value)
	case *ast.WhileStatement:
		var ret object.Object
		for {
			condition := Eval(node.Condition, env)
			if !isTrue(condition) {
				break
			}
			ret = Eval(&node.Body, env)
			if ret != nil {
				if ret.Type() == object.RETURN_VALUE_OBJ || ret.Type() == object.ERROR_OBJ {
					return ret
				}
			}
		}
		return ret
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

func evalParameters(params []ast.Expression, env *object.Environment) []object.Object {
	var ret []object.Object
	for _, p := range params {
		eval := Eval(p, env)
		if eval.Type() == object.ERROR_OBJ {
			return []object.Object{eval}
		}
		ret = append(ret, eval)
	}
	return ret
}

func evalFunction(fn object.Object, params []object.Object) object.Object {
	switch funcc := fn.(type) {
	case *object.Function:
		newEnv := expandEnv(funcc, params)
		ev := Eval(funcc.Body, newEnv)
		if retVal, ok := ev.(*object.ReturnValue); ok {
			return retVal.Value
		}
		return ev
	default:
		return &object.Error{Error: "not a func"}
	}
}

func expandEnv(fn *object.Function, params []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Params {
		env.Set(param.Value, params[i])
	}
	return env
}

func isTrue(condition object.Object) bool {
	switch condition {
	case TRUE:
		return true
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return false
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	if left.Type() != object.ARRAY_OBJ {
		return &object.Error{Error: "index expression must be applied to ARRAY object"}
	}
	if index.Type() != object.INTEGER_OBJ {
		return &object.Error{Error: "index number must be INTEGER"}
	}
	arrayObject := left.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}
