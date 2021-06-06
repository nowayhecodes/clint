package parser

import (
	"clint/ast"
	"clint/lexer"
	"fmt"
	"testing"
)

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
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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

		if !testIntegerLiteral(test, exp.RightHand, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(test *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
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

		if !testIntegerLiteral(test, exp.LeftHand, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			test.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(test, exp.RightHand, tt.rightValue) {
			return
		}

	}
}

func testIntegerLiteral(test *testing.T, intLiteral ast.Expression, value int64) bool {
	integer, ok := intLiteral.(*ast.IntegerLiteral)

	if !ok {
		test.Errorf("parameter intLiteral is not *ast.IntegerLiteral. got=%T", intLiteral)
		return false
	}

	if integer.Value != value {
		test.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		test.Errorf("integer.TokenLiteral is not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}
