package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "Illegal"
	EOF     = "EOF"

	Identifier = "Identifier"
	Integer    = "Integer"

	// Operators
	Assign        = "="
	Plus          = "+"
	PlusEqual     = "+="
	Minus         = "-"
	MinusEqual    = "-="
	Not           = "!"
	Multiply      = "*"
	MultiplyEqual = "*="
	Divide        = "/"
	DivideEqual   = "/="

	// Comparison
	Equal              = "=="
	NotEqual           = "!="
	LessThan           = "<"
	LessThanOrEqual    = "<="
	GreaterThan        = ">"
	GreaterThanOrEqual = ">="

	// Syntax
	Comma     = ","
	Semicolon = ";"

	LeftParentheses  = "("
	RightParentheses = ")"
	LeftBrace        = "{"
	RightBrace       = "}"

	// Keywords
	Function = "fn"
	Let      = "let"
	True     = "true"
	False    = "false"
	If       = "if"
	Else     = "else"
	Return   = "return"
)

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdentifier(i string) TokenType {
	if tok, ok := keywords[i]; ok {
		return tok
	}
	return Identifier
}
