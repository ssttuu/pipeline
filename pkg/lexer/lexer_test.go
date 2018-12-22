package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ssttuu/monkey/pkg/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);

!-/*5;

5 < 10 > 5;

if (5 < 100) {
	return true;
} else {
	return false;
}

10 == 10 >= 9 <= 11 != 0;

`

	tests := []struct{
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Identifier, "five"},
		{token.Assign, "="},
		{token.Integer, "5"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Identifier, "ten"},
		{token.Assign, "="},
		{token.Integer, "10"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Identifier, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LeftParentheses, "("},
		{token.Identifier, "x"},
		{token.Comma, ","},
		{token.Identifier, "y"},
		{token.RightParentheses, ")"},
		{token.LeftBrace, "{"},
		{token.Identifier, "x"},
		{token.Plus, "+"},
		{token.Identifier, "y"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Identifier, "result"},
		{token.Assign, "="},
		{token.Identifier, "add"},
		{token.LeftParentheses, "("},
		{token.Identifier, "five"},
		{token.Comma, ","},
		{token.Identifier, "ten"},
		{token.RightParentheses, ")"},
		{token.Semicolon, ";"},

		{token.Not, "!"},
		{token.Minus, "-"},
		{token.Divide, "/"},
		{token.Multiply, "*"},
		{token.Integer, "5"},
		{token.Semicolon, ";"},

		{token.Integer, "5"},
		{token.LessThan, "<"},
		{token.Integer, "10"},
		{token.GreaterThan, ">"},
		{token.Integer, "5"},
		{token.Semicolon,";"},

		{token.If, "if"},
		{token.LeftParentheses, "("},
		{token.Integer, "5"},
		{token.LessThan, "<"},
		{token.Integer, "100"},
		{token.RightParentheses, ")"},
		{token.LeftBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Else, "else"},
		{token.LeftBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},

		// 10 == 10 >= 9 <= 11 != 0;
		{token.Integer, "10"},
		{token.Equal, "=="},
		{token.Integer, "10"},
		{token.GreaterThanOrEqual, ">="},
		{token.Integer, "9"},
		{token.LessThanOrEqual, "<="},
		{token.Integer, "11"},
		{token.NotEqual, "!="},
		{token.Integer, "0"},
		{token.Semicolon, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for _, test := range tests {
		tok := l.NextToken()
		assert.Equal(t, test.expectedLiteral, tok.Literal)
		assert.Equal(t, string(test.expectedType), string(tok.Type))
	}
}
