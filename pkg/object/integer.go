package object

import "fmt"

const IntegerType ObjectType = "Integer"

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string{
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return IntegerType
}