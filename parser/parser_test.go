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
