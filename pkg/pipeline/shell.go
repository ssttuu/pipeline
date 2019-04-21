package pipeline

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"os/exec"
)

type Shell struct {
}

func (Shell) Schema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "display",
				Required: false,
			},
			{
				Name:     "script",
				Required: true,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "env",
			},
		},
	}
}

func (s *Shell) Execute(ctx *hcl.EvalContext, body hcl.Body) (map[string]cty.Value, error) {
	content, diags := body.Content(s.Schema())
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	for _, block := range content.Blocks {
		fmt.Println(" ", block.Type)
	}

	script, ok := content.Attributes["script"]
	if !ok {
		return nil, errors.New("script missing")
	}
	value, diags := script.Expr.Value(ctx)
	if diags.HasErrors() {

	}

	fmt.Println(value.AsString())

	stdout, err := exec.Command("/bin/sh", "-c", value.AsString()).Output()
	if err != nil {
		return nil, errors.Wrap(err, "running command")
	}

	output := make(map[string]cty.Value)
	output["stdout"] = cty.StringVal(string(stdout))

	return output, nil
}
