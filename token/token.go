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
	TELL    = "!"
	ASK     = "?"
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	MODULE  = "MODULE"
	CLASS   = "CLASS"
	FUN     = "FUN"
	LET     = "LET"
)

var keywords = map[string]TokenType{
	"module": MODULE,
	"class":  CLASS,
	"fun":    FUN,
	"let":    LET,
}

// LookupIdent ...
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
