package parser

import (
	"clint/ast"
	"clint/lexer"
	"fmt"
	"testing"
)

func testIdentifier(test *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		test.Errorf("exp not *ast.Identifier. got=%T", exp)
	}

	if ident.Value != value {
		test.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func TestVarStatement(test *testing.T) {
	input := `
	var x = 5;
	var y = 10;
	var foobar = 123456;
	`

	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(test, parser)

	if program == nil {
		test.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		test.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(test, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(test *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}

	test.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		test.Errorf("parser error: %q", msg)
	}
	test.FailNow()
}

func testVarStatement(test *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "var" {
		test.Errorf("stmt.TokenLiteral not 'var'. got=%q", stmt.TokenLiteral())
		return false
	}

	varStmt, ok := stmt.(*ast.VarStatement)

	if !ok {
		test.Errorf("stmt not *ast.VarStatement. got=%T", stmt)
		return false
	}

	if varStmt.Name.Value != name {
		test.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		test.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
}

// TestReturnStatement ...
func TestReturnStatement(test *testing.T) {
	input := `
	return 5;
	return 10;
	return 42;
	`

	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(test, parser)

	if len(program.Statements) != 3 {
		test.Fatalf("program.Statements does not contain 3 statements.got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			test.Errorf("Statement not *ast.returnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			test.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func TestFunctionLiteral(test *testing.T) {
	input := `fun(x, y) {x + y; }`
	lex := lexer.New(input)
	p := New(lex)
	program := p.ParseProgram()

	checkParserErrors(test, p)

	if len(program.Statements) != 1 {
		test.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		test.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		test.Fatalf("wrong numbers of function literal parameters. wanted 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(test, function.Parameters[0], "x")
	testLiteralExpression(test, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		test.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("function body is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(test, body.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(test *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fun() {};", expectedParams: []string{}},
		{input: "fun(x) {};", expectedParams: []string{"x"}},
		{input: "fun(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := New(lex)
		program := p.ParseProgram()

		checkParserErrors(test, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fun := stmt.Expression.(*ast.FunctionLiteral)

		if len(fun.Parameters) != len(tt.expectedParams) {
			test.Errorf("length parameters wrong. want=%d, got=%d\n", len(tt.expectedParams), len(fun.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(test, fun.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(test *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	lex := lexer.New(input)
	p := New(lex)
	program := p.ParseProgram()

	checkParserErrors(test, p)

	if len(program.Statements) != 1 {
		test.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		test.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(test, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		test.Fatalf("wrong number of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(test, exp.Arguments[0], 1)
	testInfixExpression(test, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(test, exp.Arguments[2], 4, "+", 5)
}

func TestIdentifierExpression(test *testing.T) {
	input := "foobar;"

	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(test, parser)

	if len(program.Statements) != 1 {
		test.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		test.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		test.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		test.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(test *testing.T) {
	input := "5;"

	lex := lexer.New(input)
	parser := New(lex)
	program := parser.ParseProgram()
	checkParserErrors(test, parser)

	if len(program.Statements) != 1 {
		test.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		test.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		test.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		test.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpression(test *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		// {"!true;", "!", true},
		// {"!false", "!", false},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		parse := New(lex)
		program := parse.ParseProgram()
		checkParserErrors(test, parse)

		if len(program.Statements) != 1 {
			test.Fatalf("program.Statements does not contain %d statemets. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			test.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			test.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(test, exp.RightHand, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(test *testing.T) {
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
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		p := New(lex)
		program := p.ParseProgram()

		checkParserErrors(test, p)

		if len(program.Statements) != 1 {
			test.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			test.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			test.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testInfixExpression(test, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}

	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.LeftHand, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.RightHand, right) {
		return false
	}

	return true
}

func TestBooleanExpression(test *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := New(lex)
		program := p.ParseProgram()
		checkParserErrors(test, p)

		if len(program.Statements) != 1 {
			test.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			test.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)

		if !ok {
			test.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}

		if boolean.Value != tt.expected {
			test.Errorf("boolean.Value not %t. got=%t", tt.expected, boolean.Value)
		}
	}
}

func TestIfExpression(test *testing.T) {
	input := `if (x < y) { x }`

	lex := lexer.New(input)
	p := New(lex)
	program := p.ParseProgram()

	checkParserErrors(test, p)

	if len(program.Statements) != 1 {
		test.Fatalf("program.Statements does not contain %d statements. got =%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		test.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(test, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		test.Errorf("consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(test, cons.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		test.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testOperatorPrecedenceParsing(test *testing.T) {
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
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
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
		lex := lexer.New(tt.input)
		p := New(lex)
		program := p.ParseProgram()
		checkParserErrors(test, p)

		actual := program.String()

		if actual != tt.expected {
			test.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
