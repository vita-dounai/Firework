package evaluator

import (
	"github.com/vita-dounai/Firework/ast"
	"github.com/vita-dounai/Firework/object"
)

func isMarcoDefinition(statement ast.Statement) (bool, *ast.MacroLiteral, string) {
	assignStatement, ok := statement.(*ast.AssignStatement)
	if !ok {
		return false, nil, ""
	}

	macroLiteral, ok := assignStatement.Value.(*ast.MacroLiteral)
	if ok {
		return true, macroLiteral, assignStatement.Name.Value
	}

	return false, nil, ""
}

func addMacroDefinition(macroLiteral *ast.MacroLiteral, env *object.Environment, name string) {
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}

	env.Set(name, macro)
}

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, statement := range program.Statements {
		ok, macroLiteral, name := isMarcoDefinition(statement)
		if ok {
			addMacroDefinition(macroLiteral, env, name)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (bool, *object.Macro) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return false, nil
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return false, nil
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return false, nil
	}

	return true, macro
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}

	for _, arg := range exp.Arguments {
		args = append(args, &object.Quote{Node: arg})
	}

	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.ExtendEnvironment(macro.Env)

	for i, parameter := range macro.Parameters {
		extended.Set(parameter.Value, args[i])
	}

	return extended
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		ok, macro := isMacroCall(callExpression, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExpression)
		extenedEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, extenedEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}
