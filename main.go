package main

import (
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	ll, err := log.ParseLevel(plugin.LogLevel)

	if err != nil {
		ll = log.InfoLevel
	}

	log.SetLevel(ll)
}

func main() {
	log.Info("--- :one: monorepo-diff")
	setupLogger()

	// pipelines := pipelinesToTrigger(plugin.Diff)
	initConfig()

	// TODO
	// get list of pipelines to trigger based on the diff
	// generate pipeline_yml
	// upload pipeline
}
