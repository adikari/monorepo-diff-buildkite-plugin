package main

import (
	"github.com/davecgh/go-spew/spew"
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

func main() {
	log.Info("--- :one: monorepo-diff")

	plugin, err := initializePlugin(env("BUILDKITE_PLUGINS", ""))

	if err != nil {
		log.Fatal(err)
	}

	setupLogger(plugin.LogLevel)
	uploadPipeline(plugin)
	spew.Dump(plugin)
}
