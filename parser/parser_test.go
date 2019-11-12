package parser

import (
	"fmt"
	"testing"

	"github.com/vita-dounai/Firework/ast"
	"github.com/vita-dounai/Firework/lexer"
)

func checkAssignStatement(t *testing.T, statement ast.Statement, name string) bool {
	assignStatement, ok := statement.(*ast.AssignStatement)
	if !ok {
		t.Errorf("statement is not *ast.AssignStatement, got=%T", statement)
		return false
	}

	if assignStatement.Name.Value != name {
		t.Errorf("assignStatement.Name.Value is not '%s', got=%q", name, assignStatement.Name.Value)
	}

	if assignStatement.Name.TokenLiteral() != name {
		t.Errorf("assignStatement.Name.TokenLiteral is not '%s', got=%q", name, assignStatement.Name)
	}

	return true
}

func checkIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d, got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d, got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func checkBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression is not *ast.Boolean, got=%T", expression)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t, got=%t", value, boolean.Value)
		return false
	}

	return true
}

func checkIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not *ast.Identifier, got=%T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %s, got=%s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral is not %s, got=%s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func checkLiteralExpression(
	t *testing.T,
	expression ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return checkIntegerLiteral(t, expression, int64(v))
	case int64:
		return checkIntegerLiteral(t, expression, v)
	case string:
		return checkIdentifier(t, expression, v)
	case bool:
		return checkBooleanLiteral(t, expression, v)
	}
	t.Errorf("unsupported type, got=%T", expression)
	return false
}

func checkInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !checkLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !checkLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestAssignStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements, got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !checkAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.AssignStatement).Value
		if !checkLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 10086;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not *ast.returnStatement, got=%T", statement)
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral() is not 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not ast.Identifier, got=%T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value is not foobar, got=%s", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral is not foobar, got=%s", identifier.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestBoolean(t *testing.T) {
	input := `
	true;
	false;
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedBoolean bool
	}{
		{true},
		{false},
	}

	for i, expected := range tests {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not ast.ExpressionStatement, got=%T", i,
				program.Statements[i])
		}

		checkBooleanLiteral(t, stmt.Expression, expected.expectedBoolean)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, world";`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral, got=%T", stmt.Expression)
	}
	if literal.Value != "Hello, world" {
		t.Errorf("literal.Value not %q, got=%q", "Hello, world", literal.Value)
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !checkLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 ** 2;", 5, "**", 2},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		checkInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	tests := []string{
		`if x < y {x}`,
	}

	for _, input := range tests {
		l := lexer.NewLexer(input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression, got=%T",
				stmt.Expression)
		}

		if !checkInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		consequenceStmt := exp.Consequence

		if len(consequenceStmt.Statements) != 1 {
			t.Errorf("consequence is not 1 statements, got=%d\n",
				len(consequenceStmt.Statements))
		}

		consequence, ok := consequenceStmt.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T",
				consequenceStmt.Statements[0])
		}
		if !checkIdentifier(t, consequence.Expression, "x") {
			return
		}

		if exp.Alternative != nil {
			t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
		}
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []string{
		`if x < y {x} else {y}`,
	}

	for _, input := range tests {
		l := lexer.NewLexer(input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression, got=%T",
				stmt.Expression)
		}

		if !checkInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		consequenceStmt := exp.Consequence
		if len(consequenceStmt.Statements) != 1 {
			t.Errorf("consequence is not 1 statements, got=%d\n",
				len(consequenceStmt.Statements))
		}

		consequence, ok := consequenceStmt.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T",
				consequenceStmt.Statements[0])
		}
		if !checkIdentifier(t, consequence.Expression, "x") {
			return
		}

		alternativeStmt := exp.Alternative
		alternative, ok := alternativeStmt.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T",
				alternativeStmt.Statements[0])
		}
		if !checkIdentifier(t, alternative.Expression, "y") {
			return
		}
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `|x, y| { x + y; }`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral, got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	checkLiteralExpression(t, function.Parameters[0], "x")
	checkLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements, got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement, got=%T",
			function.Body.Statements[0])
	}

	checkInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression, got=%T",
			stmt.Expression)
	}

	if !checkIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments, got=%d", len(exp.Arguments))
	}

	checkLiteralExpression(t, exp.Arguments[0], 1)
	checkInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	checkInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestWhileStatement(t *testing.T) {
	input := `
	while x < 10 {
		x = x + 1;
	}
	`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt is not ast.WhileStatement, got=%T",
			program.Statements[0])
	}

	if !checkInfixExpression(t, stmt.Condition, "x", "<", 10) {
		return
	}

	bodyStmt := stmt.Body

	if len(bodyStmt.Statements) != 1 {
		t.Errorf("body is not 1 statements, got=%d\n",
			len(bodyStmt.Statements))
	}

	body, ok := bodyStmt.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.AssignStatement, got=%T",
			bodyStmt.Statements[0])
	}

	if !checkLiteralExpression(t, body.Name, "x") {
		return
	}

	if !checkInfixExpression(t, body.Value, "x", "+", 1) {
		return
	}
}
