package ast

import "fmt"

type Node interface {
	fmt.Stringer
	TokenLiteral() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}
