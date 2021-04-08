package main

import "testing"

func TestUploadPipeline(t *testing.T) {
	want := "uploading pipelines"
	got := uploadPipeline()

	if want != got {
		t.Errorf(`uploadPipeline(), got %q, want "%v"`, got, want)
	}
}
