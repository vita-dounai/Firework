package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	STRING     = "STRING"

	// Operators
	ASSIGN      = "="
	PLUS        = "+"
	MINUS       = "-"
	EXCLAMATION = "!"
	ASTERISK    = "*"
	EXP         = "**"
	SLASH       = "/"
	LT          = "<"
	LTE         = "<="
	GT          = ">"
	GTE         = ">="
	EQ          = "=="
	NOT_EQ      = "!="
	VERTICAL    = "|"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	WHILE  = "WHILE"
)

var keywords = map[string]TokenType{
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENTIFIER
}
