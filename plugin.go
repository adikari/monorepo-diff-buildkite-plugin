package main

import (
	"encoding/json"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

const pluginName = "github.com/chronotc/monorepo-diff"

// HookConfig Plugin hook configuration
type HookConfig struct {
	Command string
}

// AgentConfig Plugin agent configuration
type AgentConfig struct {
	Queue string
}

// BuildConfig Plugin build configuration
type BuildConfig struct {
	Message string
	Branch  string
	Commit  string
	Env     map[string]string
}

// TriggerConfig Plugin trigger configuration
type TriggerConfig struct {
	Trigger string
	Label   string
	Build   BuildConfig
	Command string
	Async   bool
	Agents  AgentConfig
	Env     map[string]string
}

// WatchConfig Plugin watch configuration
type WatchConfig struct {
	RawPath interface{} `json:"path"`
	Paths   []string
	Path    string
	Config  TriggerConfig
}

// Plugin buildkite monorepo diff plugin structure
type Plugin struct {
	Diff          string
	Wait          bool
	LogLevel      string `json:"log_level"`
	Interpolation bool
	Hooks         []HookConfig
	Watch         []WatchConfig
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
			s.Watch[i].Path = s.Watch[i].RawPath.(string)
		case []interface{}:
			for _, v := range s.Watch[i].RawPath.([]interface{}) {
				s.Watch[i].Paths = append(s.Watch[i].Paths, v.(string))
			}
		}
		s.Watch[i].RawPath = nil
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
