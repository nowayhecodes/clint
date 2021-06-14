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

type PrefixExpression struct {
	Token     token.Token
	Operator  string
	RightHand Expression
}

func (prefix *PrefixExpression) expressionNode()      {}
func (prefix *PrefixExpression) TokenLiteral() string { return prefix.Token.Literal }
func (prefix *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prefix.Operator)
	out.WriteString(prefix.RightHand.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token     token.Token
	LeftHand  Expression
	Operator  string
	RightHand Expression
}

func (infix *InfixExpression) expressionNode()      {}
func (infix *InfixExpression) TokenLiteral() string { return infix.Token.Literal }
func (infix *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(infix.LeftHand.String())
	out.WriteString(" " + infix.Operator + " ")
	out.WriteString(infix.RightHand.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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
func (id *Identifier) String() string       { return id.Value }

// IntegerLiteral ...
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (intLiteral *IntegerLiteral) expressionNode() {}

// TokenLiteral ...
func (intLiteral *IntegerLiteral) TokenLiteral() string { return intLiteral.Token.Literal }
func (intLiteral *IntegerLiteral) String() string       { return intLiteral.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

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
