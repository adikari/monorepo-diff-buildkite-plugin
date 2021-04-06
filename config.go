package main

const pluginPrefix = "BUILDKITE_PLUGIN_MONOREPO_DIFF_"

type config struct {
	logLevel string
	diffCmd  string
}

var pluginConfig = config{
	logLevel: env(pluginPrefix+"LOG_LEVEL", "info"),
	diffCmd:  env(pluginPrefix+"DIFF", "git diff --name-only HEAD~1"),
}
