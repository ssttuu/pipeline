package parser

import (
	"fmt"
	"github.com/ssttuu/monkey/pkg/ast"
	"github.com/ssttuu/monkey/pkg/lexer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		assertIntegerLiteral(t, exp, int64(v))
	case bool:
		assertBooleanLiteral(t, exp, v)
	}
}

func assertIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	i, ok := exp.(*ast.IntegerLiteral)
	assert.True(t, ok)

	assert.Equal(t, value, i.Value)
}

func assertBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	bo, ok := exp.(*ast.Boolean)
	assert.True(t, ok)

	assert.Equal(t, value, bo.Value)
	assert.Equal(t, fmt.Sprintf("%t", value), bo.TokenLiteral())
}

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input string
		expectedIdentifier string
		expectedValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		assert.Len(t, p.errors, 0)
		assert.Len(t, program.Statements, 1)

		stmt := program.Statements[0]

		let := stmt.(*ast.LetStatement)
		assert.Equal(t, test.expectedIdentifier, let.Name.Value)
		assertLiteralExpression(t, let.Value, test.expectedValue)
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

		assertInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue)
	}
}

func assertInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) {
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
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assert.Len(t, p.errors, 0)
	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	assert.True(t, ok)

	assert.Len(t, function.Parameters, 2)

	assert.Equal(t, "x", function.Parameters[0].TokenLiteral())
	assert.Equal(t, "y", function.Parameters[1].TokenLiteral())

	assert.Len(t, function.Body.Statements, 1)

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	assertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		assert.Len(t, p.errors, 0)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		assert.Equal(t, len(test.expectedParams), len(function.Parameters))

		for i, expected := range test.expectedParams {
			assert.Equal(t, expected, function.Parameters[i].TokenLiteral())
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	assert.Len(t, p.errors, 0)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	assert.True(t, ok)

	assert.Equal(t, "add", exp.Function.TokenLiteral())
	assert.Len(t, exp.Arguments, 3)

	assert.Equal(t, "1", exp.Arguments[0].TokenLiteral())
	assertInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	assertInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}