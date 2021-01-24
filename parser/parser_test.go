package parser

import (
	"clint/ast"
	"clint/lexer"
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
