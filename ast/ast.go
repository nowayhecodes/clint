package ast

import "clint/token"

// Node ...
type Node interface {
	TokenLiteral() string
}

// Statement ...
type Statement interface {
	Node
	statementNode()
}

// Expression ...
type Expression interface {
	Node
	expressionNode()
}

// VarStatement ...
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Value Expression
}

func (vStmt *VarStatement) statementNode() {}

// TokenLiteral ...
func (vStmt *VarStatement) TokenLiteral() string { return vStmt.Token.Literal }

// ReturnStatement ...
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral ...
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// Identifier ...
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (id *Identifier) expressionNode() {}

// TokenLiteral ...
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

// Program ...
type Program struct {
	Statements []Statement
}

// TokenLiteral ...
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
