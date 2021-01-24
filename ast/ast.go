package ast

import (
	"bytes"
	"clint/token"
)

// Node ...
type Node interface {
	TokenLiteral() string
	String() string
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

// ExpressionStatement ...
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expStmt *ExpressionStatement) statementNode() {}

// TokenLiteral ...
func (expStmt *ExpressionStatement) TokenLiteral() string { return expStmt.Token.Literal }

func (expStmt *ExpressionStatement) String() string {
	if expStmt.Expression != nil {
		return expStmt.Expression.String()
	}
	return ""
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

func (vStmt *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vStmt.TokenLiteral() + " ")
	out.WriteString(vStmt.Name.String())
	out.WriteString(" = ")

	if vStmt.Value != nil {
		out.WriteString(vStmt.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// ReturnStatement ...
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral ...
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// Identifier ...
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (id *Identifier) expressionNode() {}

// TokenLiteral ...
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

func (id *Identifier) String() string { return id.Value }

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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, str := range p.Statements {
		out.WriteString(str.String())
	}
	return out.String()
}
