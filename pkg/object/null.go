package object

const NullType ObjectType = "Null"

type Null struct {

}

func (n *Null) Type() ObjectType {
	return NullType
}

func (n * Null) Inspect() string {
	return "null"
}