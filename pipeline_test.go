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
func TestUploadPipeline(t *testing.T) {
	want := "uploading pipelines"
	got := uploadPipeline()

	if want != got {
		t.Errorf(`uploadPipeline(), got %q, want "%v"`, got, want)
	}
}
