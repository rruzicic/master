package evaluator

import (
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
	case *ast.ArrayLiteral:
	case *ast.BoolLiteral:
	case *ast.StringLiteral:

	// identifier
	case *ast.IdentifierExpression:

	// expressions
	case *ast.InfixExpression:
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
