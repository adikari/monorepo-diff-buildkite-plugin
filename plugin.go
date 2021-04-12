package main

import (
	"encoding/json"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

const pluginName = "github.com/chronotc/monorepo-diff"

// Plugin buildkite monorepo diff plugin structure
type Plugin struct {
	Diff          string
	Wait          bool
	LogLevel      string `json:"log_level"`
	Interpolation bool
	Hooks         []HookConfig
	Watch         []WatchConfig
	Env           map[string]string
}

// HookConfig Plugin hook configuration
type HookConfig struct {
	Command string
}

// WatchConfig Plugin watch configuration
type WatchConfig struct {
	RawPath interface{} `json:"path"`
	Paths   []string
	Step    Step `json:"config"`
}

// Step is buildkite pipeline definition
type Step struct {
	Trigger   string            `yaml:"trigger,omitempty"`
	Label     string            `yaml:"label,omitempty"`
	Build     Build             `yaml:"build,omitempty"`
	Command   string            `yaml:"command,omitempty"`
	Agents    Agent             `yaml:"agents,omitempty"`
	Artifacts []string          `yaml:"artifacts,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Async     bool
}

// Agent is Buildkite agent definition
type Agent struct {
	Queue string `yaml:"queue,omitempty"`
}

// Build is buildkite build definition
type Build struct {
	Message string            `yaml:"message,omitempty"`
	Branch  string            `yaml:"branch,omitempty"`
	Commit  string            `yaml:"commit,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}

// UnmarshalJSON set defaults properties
func (s *Plugin) UnmarshalJSON(data []byte) error {
	type plain Plugin

	def := &plain{
		Diff:          "git diff --name-only HEAD~1",
		Wait:          false,
		LogLevel:      "info",
		Interpolation: false,
	}

	_ = json.Unmarshal(data, def)

	*s = Plugin(*def)

	for i, p := range s.Watch {
		switch p.RawPath.(type) {
		case string:
			s.Watch[i].Paths = []string{s.Watch[i].RawPath.(string)}
		case []interface{}:
			for _, v := range s.Watch[i].RawPath.([]interface{}) {
				s.Watch[i].Paths = append(s.Watch[i].Paths, v.(string))
			}
		}

		s.Watch[i].RawPath = nil

		for k, e := range s.Env {
			if s.Watch[i].Step.Env == nil {
				s.Watch[i].Step.Env = map[string]string{}
			}

			s.Watch[i].Step.Env[k] = e

			if s.Watch[i].Step.Build.Message == "" {
				continue
			}

			if s.Watch[i].Step.Build.Env == nil {
				s.Watch[i].Step.Build.Env = map[string]string{}
			}

			s.Watch[i].Step.Build.Env[k] = e
		}
	}

	return nil
}

func initializePlugin(data string) (Plugin, error) {
	var plugins []map[string]Plugin

	err := json.Unmarshal([]byte(data), &plugins)

	if err != nil {
		log.Debug(err)
		return Plugin{}, errors.New("Failed to parse plugin configuration")
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
