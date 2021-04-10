package main

import "github.com/davecgh/go-spew/spew"

func generatePipeline(pipelines []Pipeline) Pipeline {
	spew.Dump(pipelines)

	return Pipeline{}
}
