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
var Version string

func main() {
	log.Infof("--- :one: monorepo-diff %s", Version)

	plugin, err := initializePlugin(env("BUILDKITE_PLUGINS", ""))

	if err != nil {
		log.Fatal(err)
	}

	setupLogger(plugin.LogLevel)

	uploadPipeline(plugin, generatePipeline)
}
