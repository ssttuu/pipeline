package pipeline

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

var Schema = &hcl.BodySchema{}
var Plugins = map[string]Plugin{}

func init() {
	Schema.Blocks = append(Schema.Blocks, hcl.BlockHeaderSchema{
		Type: "pipeline",
		LabelNames: []string{
			"name",
		},
	})
	Plugins["pipeline"] = &Pipeline{}

	Schema.Blocks = append(Schema.Blocks, hcl.BlockHeaderSchema{
		Type: "shell",
		LabelNames: []string{
			"name",
		},
	})
	Plugins["shell"] = &Shell{}

	Schema.Blocks = append(Schema.Blocks, hcl.BlockHeaderSchema{
		Type: "stage",
		LabelNames: []string{
			"name",
		},
	})
	Plugins["stage"] = &Stage{}

	Schema.Blocks = append(Schema.Blocks, hcl.BlockHeaderSchema{
		Type: "docker-run",
		LabelNames: []string{
			"name",
		},
	})
	Plugins["docker-run"] = &DockerRun{}
}

type Pipeline struct {
}

func (Pipeline) Schema() *hcl.BodySchema {
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

func (p *Pipeline) Execute(ctx *hcl.EvalContext, body hcl.Body) (map[string]cty.Value, error) {
	content, diags := body.Content(p.Schema())
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
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
