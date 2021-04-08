package main

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	want := []string{
		"services/foo/serverless.yml",
		"services/bar/config.yml",
		"ops/bar/config.yml",
		"README.md",
	}

	got := diff("cat ./tests/mocks/diff1")

	if !reflect.DeepEqual(want, got) {
		t.Errorf(`diff(), got %q, want "%v"`, got, want)
	}
}
