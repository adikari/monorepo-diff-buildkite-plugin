package main

import "github.com/davecgh/go-spew/spew"

// Pipeline is Buildkite pipeline definition
type Pipeline struct {
	Steps []Step
}

func generatePipeline(steps []Step) Pipeline {
	spew.Dump(steps)

	return Pipeline{}
}
