package main

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

func uploadPipeline(plugin Plugin) error {
	pipelinesToTrigger := pipelinesToTrigger(
		diff(plugin.Diff),
		plugin.Watch,
	)

	generatePipeline(pipelinesToTrigger)

	return nil
}

func pipelinesToTrigger(files []string, watch []WatchConfig) []Pipeline {
	pipelines := []Pipeline{}

	for _, w := range watch {
		for _, p := range w.Paths {
			for _, f := range files {
				if strings.HasPrefix(f, p) {
					pipelines = append(pipelines, w.Config)
					break
				}
			}
		}
	}

	return dedupPipelines(pipelines)
}

func dedupPipelines(pipelines []Pipeline) []Pipeline {
	unique := []Pipeline{}
	for _, p := range pipelines {
		duplicate := false
		for _, t := range unique {
			if reflect.DeepEqual(p, t) {
				duplicate = true
				break
			}
		}

		if !duplicate {
			unique = append(unique, p)
		}
	}

	return unique
}

func diff(command string) []string {
	log.Infof("Running diff command: %s", command)

	split := strings.Split(command, " ")
	cmd, args := split[0], split[1:]

	output := executeCommand(cmd, args)

	log.Debug("Output from diff:\n" + output)

	return strings.Split(strings.TrimSpace(output), "\n")
}
