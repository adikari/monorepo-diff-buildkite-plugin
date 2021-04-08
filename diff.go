package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func diff(command string) []string {
	log.Infof("Running diff command: %s", command)

	split := strings.Split(command, " ")
	cmd, args := split[0], split[1:]

	output := executeCommand(cmd, args)

	log.Debug("Output from diff:\n" + output)

	return strings.Split(strings.TrimSpace(output), "\n")
}

func pipelinesToTrigger(diffCmd string) []string {
	changedFiles := diff(diffCmd)

	return changedFiles
}
