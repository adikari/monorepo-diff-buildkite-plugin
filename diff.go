package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func commandArgs(command string) (string, []string) {
	split := strings.Split(command, " ")

	return split[0], split[1:]
}

func diff(command string) []string {
	if command == "" {
		command = "git diff --name-only HEAD~1"
	}

	log.Infof("Running diff command: %s", command)

	output := executeCommand(commandArgs(command))

	log.Debug("Output from diff:\n" + output)

	return strings.Split(strings.TrimSpace(output), "\n")
}
