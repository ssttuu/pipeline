package object

import "fmt"

const BooleanType ObjectType = "Boolean"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BooleanType
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}