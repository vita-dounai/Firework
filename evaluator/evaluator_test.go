package evaluator

import (
	"testing"

	"github.com/vita-dounai/Firework/lexer"
	"github.com/vita-dounai/Firework/object"
	"github.com/vita-dounai/Firework/parser"
)

func checkIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func checkEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	return Eval(program)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected)
	}
}
