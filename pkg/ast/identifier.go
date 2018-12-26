package ast

import "github.com/ssttuu/monkey/pkg/token"

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
