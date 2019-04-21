package pipeline

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type DockerRun struct {
}

func (DockerRun) Schema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "display",
				Required: false,
			},
			{
				Name:     "image",
				Required: true,
			},
			{
				Name:     "command",
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

func GetAttr(attrs hcl.Attributes, key string, ctx *hcl.EvalContext) (cty.Value, error) {
	image, ok := attrs[key]
	if !ok {
		return cty.NilVal, errors.New("script missing")
	}
	value, diags := image.Expr.Value(ctx)
	if diags.HasErrors() {
		return cty.NilVal, errors.New(diags.Error())
	}

	return value, nil
}

func (s *DockerRun) Execute(ctx *hcl.EvalContext, body hcl.Body) (map[string]cty.Value, error) {
	content, diags := body.Content(s.Schema())
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	for _, block := range content.Blocks {
		fmt.Println(" ", block.Type)
	}

	image, err := GetAttr(content.Attributes, "image", ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting attribute")
	}

	commandAttr, err := GetAttr(content.Attributes, "command", ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting attribute")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "getting cwd")
	}
	command := []string{"run", "--rm", "-i", "-v", fmt.Sprintf("%s:%s", cwd, cwd), "-w", cwd, image.AsString()}
	command = append(command, strings.Split(commandAttr.AsString(), " ")...)
	cmd := exec.Command("docker", command...)

	output := make(map[string]cty.Value)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "getting stdout")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "getting stderr")
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "running command")
	}

	if bytes, err := ioutil.ReadAll(stdout); err == nil {
		output["stdout"] = cty.StringVal(string(bytes))
	}

	if bytes, err := ioutil.ReadAll(stderr); err == nil {
		output["stderr"] = cty.StringVal(string(bytes))
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.Wrap(err, "waiting for command to complete")
	}

	return output, nil
}
