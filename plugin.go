package main

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

const pluginName = "github.com/chronotc/monorepo-diff"

// Plugin buildkite monorepo diff plugin structure
type Plugin struct {
	Diff          string
	Wait          bool
	LogLevel      string `yaml:"log_level"`
	Interpolation bool
	Hooks         []struct{ Command string }
	Watch         []struct {
		Path   string
		Config struct {
			Trigger string
		}
		Label string
		Build struct {
			Message string
			Branch  string
			Commit  string
			Env     map[string]string
		}
		Command string
		Async   bool
		Agents  struct {
			Queue string
		}
		Env map[string]string
	}
}

func initializePlugin() (Plugin, error) {
	data := env("BUILDKITE_PLUGINS", "")
	var plugins []map[string]Plugin

	err := yaml.Unmarshal([]byte(data), &plugins)

	if err != nil {
		log.Debug(err)
		log.Fatal("Failed to parse plugin configuration")
	}

	for _, p := range plugins {
		for key, plugin := range p {
			if strings.HasPrefix(key, pluginName) {
				return plugin, nil
			}
		}
	}

	return Plugin{}, errors.New("Could not initialize plugin")
}
