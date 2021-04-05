package main

import "testing"

func TestUploadPipeline(t *testing.T) {
	want := "uploading pipeline"
	got := uploadPipeline()

	if want != got {
		t.Errorf(`uploadPipeline(), got %q, want "%v"`, got, want)
	}
}
