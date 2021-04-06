package main

import (
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	ll, err := log.ParseLevel(pluginConfig.logLevel)

	if err != nil {
		ll = log.InfoLevel
	}

	log.SetLevel(ll)
}

func main() {
	setupLogger()

	pipelines := pipelinesToTrigger()

	log.Info(pipelines)

	// TODO
	// get list of pipelines to trigger based on the diff
	// generate pipeline_yml
	// upload pipeline
}
