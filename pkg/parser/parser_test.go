package parser

import (
	"github.com/ssttuu/monkey/pkg/ast"
	"github.com/ssttuu/monkey/pkg/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if !assert.Len(t, p.Errors(), 0) {
		return
	}

	assert.Len(t, program.Statements, 3)


	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		assert.Equal(t, "let", stmt.TokenLiteral())
		letStmt, ok := stmt.(*ast.LetStatement)
		assert.True(t, ok)
		assert.Equal(t, test.expectedIdentifier, letStmt.Name.Value)
		assert.Equal(t, test.expectedIdentifier, letStmt.Name.TokenLiteral())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if !assert.Len(t, p.Errors(), 0) {
		return
	}

	assert.Len(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		assert.Equal(t, "return", stmt.TokenLiteral())
		_, ok := stmt.(*ast.ReturnStatement)
		assert.True(t, ok)
	}
}
