package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testGeneratePipelineWithTrigger(t *testing.T) {
	want := `
	steps:
		- trigger: foo-service
			build:
				commit: 5d93a0b58a42157201ae2ab0e7f8120ad4651489
				message: "chore: revert pipeline"
				branch: go-rewrite
		- wait
	`

	pipelines := []Pipeline{
		{
			Trigger: "foo-service-pipeline",
		},
	}

	got := generatePipeline(pipelines)

	assert.Equal(t, want, got)
}
