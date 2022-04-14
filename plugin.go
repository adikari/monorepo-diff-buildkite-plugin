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
	RawEnv        interface{} `json:"env"`
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

type Group struct {
	Label string `yaml:"group"`
	Steps []Step `yaml:"steps"`
}

// Step is buildkite pipeline definition
type Step struct {
	Group     string            `yaml:"group,omitempty"`
	Trigger   string            `yaml:"trigger,omitempty"`
	Label     string            `yaml:"label,omitempty"`
	Build     Build             `yaml:"build,omitempty"`
	Command   string            `yaml:"command,omitempty"`
	Agents    map[string]string `yaml:"agents,omitempty"`
	Artifacts []string          `yaml:"artifacts,omitempty"`
	RawEnv    interface{}       `json:"env" yaml:",omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Async     bool              `yaml:"async,omitempty"`
}

// Build is buildkite build definition
type Build struct {
	Message string            `yaml:"message,omitempty"`
	Branch  string            `yaml:"branch,omitempty"`
	Commit  string            `yaml:"commit,omitempty"`
	RawEnv  interface{}       `json:"env" yaml:",omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
}

func (s Step) MarshalYAML() (interface{}, error) {
	if s.Group == "" {
		type Alias Step
		return (Alias)(s), nil
	}

	label := s.Group
	s.Group = ""
	return Group{Label: label, Steps: []Step{s}}, nil
}

func initializePlugin(data string) (Plugin, error) {
	var plugins []map[string]Plugin

	err := json.Unmarshal([]byte(data), &plugins)

	if err != nil {
		log.Debug(err)
		return Plugin{}, errors.New("failed to parse plugin configuration")
	}

	for _, p := range plugins {
		for key, plugin := range p {
			if strings.HasPrefix(key, pluginName) {
				return plugin, nil
			}
		}
	}

	return Plugin{}, errors.New("could not initialize plugin")
}

// UnmarshalJSON set defaults properties
func (plugin *Plugin) UnmarshalJSON(data []byte) error {
	type plain Plugin

	def := &plain{
		Diff:          "git diff --name-only HEAD~1",
		Wait:          false,
		LogLevel:      "info",
		Interpolation: true,
	}

	_ = json.Unmarshal(data, def)

	*plugin = Plugin(*def)

	plugin.Env = parseEnv(plugin.RawEnv)
	plugin.RawEnv = nil

	// Path can be string or an array of strings,
	// handle both cases and create an array of paths.
	for i, p := range plugin.Watch {
		switch p.RawPath.(type) {
		case string:
			plugin.Watch[i].Paths = []string{plugin.Watch[i].RawPath.(string)}
		case []interface{}:
			for _, v := range plugin.Watch[i].RawPath.([]interface{}) {
				plugin.Watch[i].Paths = append(plugin.Watch[i].Paths, v.(string))
			}
		}

		if plugin.Watch[i].Step.Trigger != "" {
			setBuild(&plugin.Watch[i].Step.Build)
		}

		appendEnv(&plugin.Watch[i], plugin.Env)

		p.RawPath = nil
	}

	return nil
}

func setBuild(build *Build) {
	if build.Message == "" {
		build.Message = env("BUILDKITE_MESSAGE", "")
	}

	if build.Branch == "" {
		build.Branch = env("BUILDKITE_BRANCH", "")
	}

	if build.Commit == "" {
		build.Commit = env("BUILDKITE_COMMIT", "")
	}
}

// appends top level env to Step.Env and Step.Build.Env
func appendEnv(watch *WatchConfig, env map[string]string) {
	watch.Step.Env = parseEnv(watch.Step.RawEnv)
	watch.Step.Build.Env = parseEnv(watch.Step.Build.RawEnv)

	for key, value := range env {
		if watch.Step.Command != "" {
			if watch.Step.Env == nil {
				watch.Step.Env = make(map[string]string)
			}

			watch.Step.Env[key] = value
			continue
		}
		if watch.Step.Trigger != "" {
			if watch.Step.Build.Env == nil {
				watch.Step.Build.Env = make(map[string]string)
			}

			watch.Step.Build.Env[key] = value
			continue
		}
	}

	watch.Step.RawEnv = nil
	watch.Step.Build.RawEnv = nil
	watch.RawPath = nil
}

// parse env in format from env=env-value to map[env] = env-value
func parseEnv(raw interface{}) map[string]string {
	if raw == nil {
		return nil
	}

	result := make(map[string]string)
	for _, v := range raw.([]interface{}) {
		split := strings.Split(v.(string), "=")
		key, value := strings.TrimSpace(split[0]), split[1:]

		// only key exists. set value from env
		if len(key) > 0 && len(value) == 0 {
			result[key] = env(key, "")
		}

		if len(value) > 0 {
			result[key] = strings.TrimSpace(value[0])
		}
	}

	return result
}
