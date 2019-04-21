package main

import (
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/ssttuu/pipeline/pkg/pipeline"
	"log"
)

func main() {
	var runner pipeline.Runner = &pipeline.Local{}

	p := hclparse.NewParser()
	f, _ := p.ParseHCLFile("pr.pipeline")

	ctx := pipeline.NewContext()
	if err := runner.Run(ctx, f.Body); err != nil {
		log.Fatalln(err)
	}

}
