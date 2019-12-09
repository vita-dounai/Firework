package evaluator

import (
	"github.com/vita-dounai/Firework/ast"
	"github.com/vita-dounai/Firework/object"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	if name, ok := callExpression.Function.(*ast.Identifier); ok {
		return name.Value == "unquote"
	}

	return false
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(callExpression.Arguments) != 1 {
			return node
		}

		unquoted := Eval(callExpression.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		return &ast.IntegerLiteral{Value: obj.Value}
	case *object.Boolean:
		return &ast.Boolean{Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
