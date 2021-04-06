package main

const pluginPrefix = "BUILDKITE_PLUGIN_MONOREPO_DIFF_"

type config struct {
	logLevel string
}

var pluginConfig = config{
	logLevel: env(pluginPrefix+"LOG_LEVEL", "info"),
}
