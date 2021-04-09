package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginWithEmptyParameter(t *testing.T) {
	_, err := initializePlugin("[]")

	assert.EqualError(t, err, "Could not initialize plugin")
}

func TestPluginWithInvalidParameter(t *testing.T) {
	_, err := initializePlugin("invalid")

	assert.EqualError(t, err, "Failed to parse plugin configuration")
}

func TestPluginShouldHaveDefaultValues(t *testing.T) {
	param := `[{
		"github.com/chronotc/monorepo-diff-buildkite-plugin#commit": {}
	}]`

	got, _ := initializePlugin(param)

	expected := Plugin{
		Diff:          "git diff --name-only HEAD~1",
		Wait:          false,
		LogLevel:      "info",
		Interpolation: false,
	}

	assert.Equal(t, expected, got)
}

func TestPluginWithValidParameter(t *testing.T) {
	param := ""
	got, _ := initializePlugin(param)
	expected := Plugin{}

	assert.Equal(t, expected, got)
}

func TestPluginShouldUnmarshallCorrectly(t *testing.T) {
	param := `[{
		"github.com/chronotc/monorepo-diff-buildkite-plugin#commit": {
			"diff": "cat ./hello.txt",
			"wait": true,
			"log_level": "debug",
			"interpolation": true,
			"hooks": [
				{ "command": "some-hook-command" },
				{ "command": "another-hook-command" }
			],
			"watch": [
				{
					"path": "watch-path-1",
					"config": {
						"trigger": "service-2"
					}
				},
				{
					"path": [
						"watch-path-1",
						"watch-path-2"
					],
					"config": {
						"trigger": "service-1",
						"label": "hello",
						"build": {
							"message": "build message",
							"branch": "current branch",
							"commit": "commit-hash",
							"env": {
								"foo": "bar",
								"bar": "foo"
							}
						},
						"async": true,
						"agents": {
							"queue": "queue-1"
						},
						"env": {
							"foo": "bar"
						}
					}
				}
			]
		}
	}]`

	got, _ := initializePlugin(param)

	expected := Plugin{
		Diff:          "cat ./hello.txt",
		Wait:          true,
		LogLevel:      "debug",
		Interpolation: true,
		Hooks: []HookConfig{
			{Command: "some-hook-command"},
			{Command: "another-hook-command"},
		},
		Watch: []WatchConfig{
			{Path: "watch-path-1", Config: TriggerConfig{Trigger: "service-2"}},
			{
				Paths: []string{"watch-path-1", "watch-path-2"},
				Config: TriggerConfig{
					Trigger: "service-1",
					Label:   "hello",
					Build: BuildConfig{
						Message: "build message",
						Branch:  "current branch",
						Commit:  "commit-hash",
						Env: map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
					},
					Async:  true,
					Agents: AgentConfig{Queue: "queue-1"},
					Env: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
	}

	assert.Equal(t, expected, got)
}
