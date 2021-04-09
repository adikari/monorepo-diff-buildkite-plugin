package main

const pluginName = "chronotc/monorepo-diff"
const pluginPrefix = "BUILDKITE_PLUGIN_MONOREPO_DIFF_"

// Plugin buildkite monorepo diff plugin structure
type Plugin struct {
	Diff          string
	Wait          bool
	LogLevel      string `yaml:"log_level"`
	Interpolation bool
	Hooks         []struct{ Command string }
	Watch         []struct {
		Path   string
		Config struct {
			Trigger string
		}
		Label string
		Build struct {
			Message string
			Branch  string
			Commit  string
			Env     map[string]string
		}
		Command string
		Async   bool
		Agents  struct {
			Queue string
		}
		Env map[string]string
	}
}

var plugin = Plugin{
	Diff:     readProperty("diff", "git diff --name-only HEAD~1"),
	Wait:     readBool("wait", "false"),
	LogLevel: readProperty("log_level", "info"),
}

func initConfig() {}
