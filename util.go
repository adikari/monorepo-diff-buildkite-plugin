package main

import (
	"bytes"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func executeCommand(command string, args []string) string {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Debugf(
			"\ncommand = '%s', \nargs = '%s', \nerror = '%s'",
			command, args, stderr.String(),
		)
		log.Fatalf("'%s' command failed.", command)
	}

	return out.String()
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
