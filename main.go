package main

import (
	"os"

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

	plugin, err := initializePlugin()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	setupLogger(plugin.LogLevel)

	// pipelines := pipelinesToTrigger(plugin.Diff)

	// TODO
	// get list of pipelines to trigger based on the diff
	// generate pipeline_yml
	// upload pipeline
}
