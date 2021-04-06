package main

import (
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})

	ll, err := log.ParseLevel(PluginConfig.logLevel)

	if err != nil {
		ll = log.InfoLevel
	}

	log.SetLevel(ll)
}

func main() {
	setupLogger()

	changed := diff("")

	log.Debug("debug message")
	log.Info(changed)
	// perform diff
	// get list of pipelines to trigger based on the diff
	// generate pipeline_yml
	// upload pipeline
}
