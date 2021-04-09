package main

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func executeCommand(command string, args []string) string {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Debugf(
			"\ncommand = '%s', \nargs = '%s', \nerror = '%s'",
			command, args, err.Error(),
		)
		log.Fatalf("'%s' command failed.", command)
	}

	return string(stdout)
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
