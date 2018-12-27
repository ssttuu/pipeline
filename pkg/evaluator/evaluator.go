package evaluator

import (
	"github.com/ssttuu/monkey/pkg/ast"
	"github.com/ssttuu/monkey/pkg/object"
	"github.com/ssttuu/monkey/pkg/token"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.Not:
		return evalNotOperatorExpression(right)
	case token.Minus:
		return evalMinusPrefixOperatorExpression(right)
	default:
		return Null
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == token.Equal:
		return nativeBoolToBooleanObject(left == right)
	case operator == token.NotEqual:
		return nativeBoolToBooleanObject(left != right)
	default:
		return Null
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.Plus:
		return &object.Integer{Value: leftVal + rightVal}
	case token.Minus:
		return &object.Integer{Value: leftVal - rightVal}
	case token.Multiply:
		return &object.Integer{Value: leftVal * rightVal}
	case token.Divide:
		return &object.Integer{Value: leftVal / rightVal}
	case token.LessThan:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case token.LessThanOrEqual:
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case token.GreaterThan:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case token.GreaterThanOrEqual:
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case token.Equal:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case token.NotEqual:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return Null
	}
}

func evalNotOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerType {
		return Null
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}