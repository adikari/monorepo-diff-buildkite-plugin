package main

import (
	"bytes"
	"fmt"
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

		return "", fmt.Errorf("command `%s` failed: %v", command, err)
	}

	return out.String(), nil
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func isString(val interface{}) (string, bool) {
	if val == nil {
		return "", false
	}

	switch val.(type) {
	case string:
		return val.(string), true
	}

	return "", false
}
