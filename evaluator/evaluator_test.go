package evaluator

import (
	"testing"

	"github.com/vita-dounai/Firework/lexer"
	"github.com/vita-dounai/Firework/object"
	"github.com/vita-dounai/Firework/parser"
)

func checkNullObject(t *testing.T, obj object.Object) bool {
	return obj == NULL
}

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

func checkBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func checkEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5 ** 2", 25},
		{"5 + 5 ** 2", 30},
		{"(5 + 5) ** 2", 100},
		{"5 * 5 ** 2", 125},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		checkBooleanObject(t, evaluated, tt.expected)
	}
}

func TestExclamationOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		checkBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			checkIntegerObject(t, evaluated, int64(integer))
		} else {
			checkNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`if(10 > 1) {
				if(10 > 1) {
					return 10;
				}

				return 1;
			}`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
					return 1;
				}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)

		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if err.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected=%q, got=%q",
				tt.expectedMessage, err.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		checkIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "|x| {x + 2;};"

	evaluated := checkEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters, Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = |x| { x; }; identity(5);", 5},
		{"let identity = |x| { return x; }; identity(5);", 5},
		{"let double = |x| { x * 2; }; double(5);", 10},
		{"let add = |x, y| { x + y; }; add(5, 5);", 10},
		{"let add = |x, y| { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"|x| { x; }(5)", 5},
	}
	for _, tt := range tests {
		checkIntegerObject(t, checkEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, world"`

	evaluated := checkEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello, world" {
		t.Errorf("String has wrong value, got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + ", " + "world"`
	evaluated := checkEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello, world" {
		t.Errorf("String has wrong value, got=%q", str.Value)
	}
}

func TestStringComp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"a" == "a"`, true},
		{`"a" == "b"`, false},
		{`"a" != "b"`, true},
		{`"a" < "b"`, true},
		{`"a" > "b"`, false},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)
		result, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("object is not Boolean, got=%T (%+v)", evaluated, evaluated)
		}

		if result.Value != tt.expected {
			t.Fatalf("result has wrong value, got=%t, want=%t",
				result.Value, tt.expected)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := checkEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			checkIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestWhileStatement(t *testing.T) {
	input := `
	let x = 1;
	while x < 10 {
		let x = x + 1;
	}
	x;
	`

	evaluated := checkEval(input)
	if !checkIntegerObject(t, evaluated, 10) {
		return
	}
}
