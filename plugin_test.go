package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginWithEmptyParameter(t *testing.T) {
	_, err := initializePlugin("")

	assert.EqualError(t, err, "Could not initialize plugin")
}

func TestPluginWithInvalidParameter(t *testing.T) {
	_, err := initializePlugin("invalid")

	assert.EqualError(t, err, "Failed to parse plugin configuration")
}

func TestPluginShouldHaveDefaultValues(t *testing.T) {
	param := ""
	got, _ := initializePlugin(param)
	expected := Plugin{
		Diff: "git diff --name-only HEAD~1",
	}

	assert.Equal(t, expected, got)
}

func TestPluginWithValidParameter(t *testing.T) {
	param := ""
	got, _ := initializePlugin(param)
	expected := Plugin{
		Diff: "",
	}

	assert.Equal(t, expected, got)
}
