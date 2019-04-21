package pipeline

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type Runner interface {
	Run(ctx *hcl.EvalContext, body hcl.Body) error
}

type Local struct {
}

func (l *Local) Run(ctx *hcl.EvalContext, body hcl.Body) error {
	content, _ := body.Content(Schema)
	for _, block := range content.Blocks {
		fmt.Println(block.Type, block.TypeRange)
		blockCtx := NewChild(ctx)
		output, err := Plugins[block.Type].Execute(blockCtx, block.Body)
		if err != nil {
			return errors.Wrapf(err, "executing plugin: %v", block.Type)
		}
		ctx.Variables[block.Labels[0]] = cty.ObjectVal(output)
	}
	return nil
}
