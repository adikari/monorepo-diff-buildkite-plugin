package main

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

func uploadPipeline(plugin Plugin) error {
	steps := stepsToTrigger(
		diff(plugin.Diff),
		plugin.Watch,
	)

	generatePipeline(steps)

	return nil
}

func stepsToTrigger(files []string, watch []WatchConfig) []Step {
	steps := []Step{}

	for _, w := range watch {
		for _, p := range w.Paths {
			for _, f := range files {
				if strings.HasPrefix(f, p) {
					steps = append(steps, w.Step)
					break
				}
			}
		}
	}

	return dedupSteps(steps)
}

func dedupSteps(steps []Step) []Step {
	unique := []Step{}
	for _, p := range steps {
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
