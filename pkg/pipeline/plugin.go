package pipeline

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Plugin interface {
	Schema() *hcl.BodySchema
	Execute(ctx *hcl.EvalContext, data hcl.Body) (map[string]cty.Value, error)
}
