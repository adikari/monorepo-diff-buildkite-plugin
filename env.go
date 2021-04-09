package main

import (
	"os"
	"strconv"
	"strings"
)

func readProperty(prop string, def string) string {
	return env(pluginPrefix+strings.ToUpper(prop), def)
}

func readBool(prop string, def string) bool {
	v, _ := strconv.ParseBool(readProperty(prop, def))
	return v
}

func env(key, fb string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fb
}
