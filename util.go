package main

import (
	"log"
	"os/exec"
)

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func executeCommand(command string, args []string) string {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(stdout)
}
