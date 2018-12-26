package ast

import (
	"bytes"
	"github.com/ssttuu/monkey/pkg/token"
	"strings"
)

type FunctionLiteral struct {
	Token token.Token
	Parameters []*Identifier
	Body *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())
	return out.String()
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
