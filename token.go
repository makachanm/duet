package main

// TokenType is a string representing the type of a token.
type TokenType string

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	FLOAT  = "FLOAT"  // 3.14
	STRING = "STRING" // "hello world"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	LT       = "<"
	GT       = ">"
	LE       = "<="
	GE       = ">="
	EQ       = "=="
	NOT_EQ   = "!="
	ARROW    = "->"
	PIPELINE = "|>"

	// Delimiters
	COMMA    = ","
	COLON    = ":"
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	PROC    = "PROC"
	CONS    = "CONS"
	SUPP    = "SUPP"
	IF      = "IF"
	THEN    = "THEN"
	ELSE    = "ELSE"
	FOR     = "FOR"
	IN      = "IN"
	FAIL    = "FAIL"
	MATCH   = "MATCH"
	IS      = "IS"
	DEFAULT = "DEFAULT"
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	NIL     = "NIL"
)

var keywords = map[string]TokenType{
	"proc":    PROC,
	"cons":    CONS,
	"supp":    SUPP,
	"if":      IF,
	"then":    THEN,
	"else":    ELSE,
	"for":     FOR,
	"in":      IN,
	"match":   MATCH,
	"is":      IS,
	"default": DEFAULT,
	"fail":    FAIL,
	"true":    TRUE,
	"false":   FALSE,
	"nil":     NIL,
}

// LookupIdent checks the keywords table to see whether the given identifier is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
