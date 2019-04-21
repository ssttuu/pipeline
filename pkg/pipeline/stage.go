package pipeline

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)


type Stage struct {
}

func (Stage) Schema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "display",
				Required: false,
			},
		},
		Blocks: append([]hcl.BlockHeaderSchema{
			{
				Type:       "env",
				LabelNames: []string{},
			},
		}, Schema.Blocks...),
	}
}

func (p *Stage) Execute(ctx *hcl.EvalContext, body hcl.Body) (map[string]cty.Value, error) {
	content, diags := body.Content(p.Schema())
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	for key, attr := range content.Attributes {
		value, err := attr.Expr.Value(ctx)
		if !err.HasErrors() {
			fmt.Println(" ", key, value.AsString())
		}
	}

	for _, block := range content.Blocks {
		name := block.Labels[0]
		fmt.Println(block.Type, name)
		stepCtx := NewChild(ctx)
		output, err := Plugins[block.Type].Execute(stepCtx, block.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "executing plugin: %v", block.Type)
		}
		ctx.Variables[name] = cty.ObjectVal(output)
	}

	return ctx.Variables, nil
}
