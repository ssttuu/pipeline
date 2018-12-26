package parser

import (
	"fmt"
	"github.com/ssttuu/monkey/pkg/ast"
	"github.com/ssttuu/monkey/pkg/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case bool:
		testBooleanLiteral(t, exp, v)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	bo, ok := exp.(*ast.Boolean)
	assert.True(t, ok)

	assert.Equal(t, value, bo.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), bo.TokenLiteral())
}

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

func TestBooleanExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{`true;`, true},
		{`false;`, false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		exp, ok := stmt.Expression.(*ast.Boolean)
		assert.True(t, ok)
		assert.Equal(t, test.expected, exp.Value)
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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, test := range infixTests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		assert.Len(t, program.Statements, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		testInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue)
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) {
	infix, ok := exp.(*ast.InfixExpression)
	assert.True(t, ok)

	testNode(t, left, infix.Left)
	assert.Equal(t, operator, infix.Operator)
	testNode(t, right, infix.Right)
}

func testNode(t *testing.T, value interface{}, expected interface{}) {
	switch v := value.(type) {
	case int:
		assert.Equal(t, int64(v), expected.(*ast.IntegerLiteral).Value)
	case bool:
		assert.Equal(t, v, expected.(*ast.Boolean).Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!true;", "!", true},
		{"!false;", "!", false},
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"5 >= 3 == 3 <= 5",
			"((5 >= 3) == (3 <= 5))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()
		assert.Len(t, p.Errors(), 0)

		assert.Equal(t, test.expected, program.String())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assert.Len(t, p.errors, 0)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	assert.Len(t, exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	assert.Equal(t, "x", consequence.Expression.TokenLiteral())
	assert.Nil(t, exp.Alternative)
}

func testIdentifier(t *testing.T, exp ast.Expression, v string) {


}
