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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	identifier, ok := stmt.Expression.(*ast.Identifier)
	assert.True(t, ok)

	assert.Equal(t, "foobar", identifier.Value)
	assert.Equal(t, "foobar", identifier.TokenLiteral())
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	assert.True(t, ok)

	assert.Equal(t, int64(5), literal.Value)
	assert.Equal(t, "5", literal.TokenLiteral())
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		value int64
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
	}

	for _, test := range prefixTests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		assert.Len(t, program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok)
		assert.Equal(t, test.operator, exp.Operator)

	}
}
