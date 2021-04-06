package main

const pluginPrefix = "BUILDKITE_PLUGIN_MONOREPO_DIFF_"

type config struct {
	logLevel string
}

// PluginConfig is map of all configs for plugin
var PluginConfig = config{
	logLevel: env(pluginPrefix+"LOG_LEVEL", "info"),
}
