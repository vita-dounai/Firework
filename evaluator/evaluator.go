package evaluator

import (
	"fmt"
	"strings"

	"github.com/vita-dounai/Firework/ast"
	"github.com/vita-dounai/Firework/object"
)

var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	BREAK    = &object.LoopControl{ControlType: object.BREAK}
	CONTINUE = &object.LoopControl{ControlType: object.CONTINUE}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	} else {
		return FALSE
	}
}

func evalExclamationOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		if right.Value != 0 {
			return FALSE
		} else {
			return TRUE
		}
	case *object.Boolean:
		return nativeBoolToBooleanObject(!right.Value)
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalExclamationOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "**":
		result := int64(1)
		for i := rightValue; i > 0; i >>= 1 {
			if i&1 != 0 {
				result *= leftValue
			}
			leftValue *= leftValue
		}

		return &object.Integer{Value: result}
	case "%":
		return &object.Integer{Value: leftValue % rightValue}
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case ">=":
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "<=":
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator: %s %s  %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "<":
		return nativeBoolToBooleanObject(strings.Compare(leftValue, rightValue) < 0)
	case ">":
		return nativeBoolToBooleanObject(strings.Compare(leftValue, rightValue) > 0)
	case "==":
		return nativeBoolToBooleanObject(strings.Compare(leftValue, rightValue) == 0)
	case "!=":
		return nativeBoolToBooleanObject(strings.Compare(leftValue, rightValue) != 0)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		switch {
		case left.Type() != right.Type():
			return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
		default:
			return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
		}
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	extendedEnv := object.ExtendEnvironment(env)
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, extendedEnv)

		if result != nil {
			switch result.Type() {
			case object.RETURN_VALUE_OBJ:
				fallthrough
			case object.ERROR_OBJ:
				fallthrough
			case object.BREAK:
				fallthrough
			case object.CONTINUE:
				return result
			}
		}
	}

	return result
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func isReturn(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.RETURN_VALUE_OBJ
	}

	return false
}

func isBreak(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.BREAK
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(node.Value); ok {
		return value
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("Identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	extendedEnv := object.ExtendEnvironment(fn.Env)

	for idx, param := range fn.Parameters {
		extendedEnv.Set(param.Value, args[idx])
	}

	return extendedEnv
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return newError("Not a function: %s", fn.Type())
	}
}

func evalIndexExpression(leftObject, indexObject object.Object) object.Object {
	switch left := leftObject.(type) {
	case *object.Array:
		index, ok := indexObject.(*object.Integer)
		if !ok {
			return newError("Subscript not support: %s", index.Type())
		}

		subscript := index.Value

		if subscript < 0 || subscript >= int64(len(left.Elements)) {
			return NULL
		}

		return left.Elements[subscript]
	case *object.Map:
		index, ok := indexObject.(object.Hashable)

		if !ok {
			return newError("unusable as map key: %s", indexObject.Type())
		}

		hashKey := index.Hash()
		pair, ok := left.Pairs[hashKey]
		if !ok {
			return NULL
		}

		return pair.Value
	}
	return NULL
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		returnValue := Eval(node.ReturnValue, env)
		if isError(returnValue) {
			return returnValue
		}
		return &object.ReturnValue{Value: returnValue}
	case *ast.AssignStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body
		return &object.Function{Parameters: parameters, Body: body, Env: env}
	case *ast.CallExpression:
		if name, ok := node.Function.(*ast.Identifier); ok && name.Value == "quote" {
			return quote(node.Arguments[0], env)
		}

		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.WhileStatement:
		for true {
			condition := Eval(node.Condition, env)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				break
			}

			body := Eval(node.Body, env)

			if isError(body) || isReturn(body) {
				return body
			}

			if isBreak(body) {
				break
			}
		}
	case *ast.BreakStatement:
		return BREAK
	case *ast.ContinueStatement:
		return CONTINUE
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		if left.Type() != object.ARRAY_OBJ && left.Type() != object.MAP_OBJ {
			return newError("Index operator not support: %s", left.Type())
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.MapLiteral:
		pairs := make(map[object.HashKey]object.MapPair)

		for keyNode, valueNode := range node.Pairs {
			key := Eval(keyNode, env)
			if isError(key) {
				return key
			}

			hashableKeyObject, ok := key.(object.Hashable)
			if !ok {
				return newError("unusable as map key: %s", key.Type())
			}

			value := Eval(valueNode, env)
			if isError(value) {
				return value
			}

			hashKey := hashableKeyObject.Hash()
			pairs[hashKey] = object.MapPair{Key: key, Value: value}
		}

		return &object.Map{Pairs: pairs}
	}

	return nil
}
