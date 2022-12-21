package main

import (
	log "github.com/sirupsen/logrus"
)

func setupLogger(logLevel string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	ll, err := log.ParseLevel(logLevel)

	if err != nil {
		ll = log.InfoLevel
	}

	log.SetLevel(ll)
}

// Version of plugin
var version string = "dev"

func main() {
	log.Infof("--- running monorepo-diff-buildkite-plugin %s", version)

	plugins := env("BUILDKITE_PLUGINS", "")

	log.Debugf("received plugin: \n%v", plugins)

	plugin, err := initializePlugin(plugins)

	if err != nil {
		log.Debug(err)
		log.Fatal(err)
	}

	setupLogger(plugin.LogLevel)

	if _, _, err = uploadPipeline(plugin, generatePipeline); err != nil {
		log.Fatalf("+++ failed to upload pipeline: %v", err)
	}
}
