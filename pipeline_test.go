package main

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestDiff(t *testing.T) {
	want := []string{
		"services/foo/serverless.yml",
		"services/bar/config.yml",
		"ops/bar/config.yml",
		"README.md",
	}

	got := diff("cat ./tests/mocks/diff1")

	assert.Equal(t, want, got)
}

func TestPipelinesToTriggerGetsListOfPipelines(t *testing.T) {
	want := []string{"service-1", "service-2", "service-4"}

	watch := []WatchConfig{
		{
			Paths: []string{"watch-path-1"},
			Step:  Step{Trigger: "service-1"},
		},
		{
			Paths: []string{"watch-path-2/", "watch-path-3/", "watch-path-4"},
			Step:  Step{Trigger: "service-2"},
		},
		{
			Paths: []string{"watch-path-5"},
			Step:  Step{Trigger: "service-3"},
		},
		{
			Paths: []string{"watch-path-2"},
			Step:  Step{Trigger: "service-4"},
		},
	}

	changedFiles := []string{
		"watch-path-1/text.txt",
		"watch-path-2/.gitignore",
		"watch-path-2/src/index.go",
		"watch-path-4/test/index_test.go",
	}

	pipelines := stepsToTrigger(changedFiles, watch)
	var got []string

	for _, v := range pipelines {
		got = append(got, v.Trigger)
	}

	assert.Equal(t, want, got)
}

// func TestUploadPipeline(t *testing.T) {
// 	want := "uploading pipelines"
// 	got := uploadPipeline()

// 	if want != got {
// 		t.Errorf(`uploadPipeline(), got %q, want "%v"`, got, want)
// 	}
// }
