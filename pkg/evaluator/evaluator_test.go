package evaluator

import (
	"github.com/ssttuu/monkey/pkg/lexer"
	"github.com/ssttuu/monkey/pkg/object"
	"github.com/ssttuu/monkey/pkg/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{"5", 5},
		{"-5", -5},
		{"10", 10},
		{"-10", -10},
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		assertObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	return Eval(program)
}

func assertObject(t *testing.T, obj object.Object, expected interface{}) {
	switch v := expected.(type) {
	case int:
		assertIntegerObject(t, obj, int64(v))
	case bool:
		assertBooleanObject(t, obj, v)
	}
}

func assertBooleanObject(t *testing.T, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	assert.True(t, ok)

	assert.Equal(t, expected, result.Value)
}

func assertIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	assert.True(t, ok)

	assert.Equal(t, expected, result.Value)
}
