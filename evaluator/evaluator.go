package evaluator

import (
	"interpreter/ast"
	"interpreter/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node, env)
	}
	return nil
}

func evalExpressionStatement(node *ast.ExpressionStatement, env *object.Environment) object.Object {

	if i, ok := node.Expression.(*ast.IntegerLiteral); ok {
		return &object.Integer{Value: i.Value}
	}
	return nil
}
