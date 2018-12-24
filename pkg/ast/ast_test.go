package ast

import (
	"github.com/ssttuu/monkey/pkg/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	p := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal(t, `let myVar = anotherVar;`, p.String())
}
