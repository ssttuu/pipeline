package pipeline

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func NewContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
		Functions: make(map[string]function.Function),
	}
}

func NewChild(ctx *hcl.EvalContext) *hcl.EvalContext {
	child := ctx.NewChild()
	child.Variables = make(map[string]cty.Value)
	child.Functions = make(map[string]function.Function)
	return child
}
