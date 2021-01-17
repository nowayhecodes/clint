package token

// TokenType ...
type TokenType string

// Token ...
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"
	INT     = "INT"
	FLOAT   = "FLOAT"
	STR     = "STRING"
	ASSIGN  = "="
	PLUS    = "+"
	MINUS   = "-"
	MULT    = "*"
	DIV     = "/"
	MOD     = "%"
	POW     = "^"
	COMMA   = ","
	COLON   = ":"
	SEMI    = ";"
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	MODULE  = "MODULE"
	CLASS   = "CLASS"
	FUN     = "FUN"
	LET     = "LET"
)
