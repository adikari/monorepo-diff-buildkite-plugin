package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func executeCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Debugf(
			"\ncommand = '%s', \nargs = '%s', \nerror = '%s'",
			command, args, stderr.String(),
		)

		return "", errors.New("command failed")
	}

	return out.String(), nil
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
