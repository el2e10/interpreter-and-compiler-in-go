package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF" // Represent end of the file

	IDENT = "IDENT" // Identifiers like function name, variable names
	INT   = "INT"   // Integers

	ASSIGN   = "="
	EQ       = "=="
	NOT_EQ   = "!="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	COMMA     = ","
	SEMICOLON = ";"

	LT = "<"
	GT = ">"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords -> language specific words like type, func in Go
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"

	STRING = "STRING" 
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func LookupIdentifier(token string) TokenType {
	if value, ok := keywords[token]; ok {
		return value
	}
	return IDENT
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
