package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// WaitStep represents a Buildkite Wait Step
// https://buildkite.com/docs/pipelines/wait-step
// We can't use Step here since the value for Wait is always nil
// regardless of whether or not we want to include the key.
type WaitStep struct{}

func (WaitStep) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"wait": nil,
	}, nil
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

func (n PluginNotify) MarshalYAML() (interface{}, error) {
	return n, nil
}

// PipelineGenerator generates pipeline file
type PipelineGenerator func(steps []Step, plugin Plugin) (*os.File, error)

func uploadPipeline(plugin Plugin, generatePipeline PipelineGenerator) (string, []string, error) {
	diffOutput, err := diff(plugin.Diff)
	if err != nil {
		log.Fatal(err)
		return "", []string{}, err
	}

	if len(diffOutput) < 1 {
		log.Info("No changes detected. Skipping pipeline upload.")
		return "", []string{}, nil
	}

	log.Debug("Output from diff: \n" + strings.Join(diffOutput, "\n"))

	steps, err := stepsToTrigger(diffOutput, plugin.Watch)
	if err != nil {
		return "", []string{}, err
	}

	pipeline, err := generatePipeline(steps, plugin)
	defer os.Remove(pipeline.Name())

	if err != nil {
		log.Error(err)
		return "", []string{}, err
	}

	cmd := "buildkite-agent"
	args := []string{"pipeline", "upload", pipeline.Name()}

	if !plugin.Interpolation {
		args = append(args, "--no-interpolation")
	}

	_, err = executeCommand("buildkite-agent", args)

	return cmd, args, err
}

func diff(command string) ([]string, error) {
	log.Infof("Running diff command: %s", command)

	output, err := executeCommand(
		env("SHELL", "bash"),
		[]string{"-c", strings.Replace(command, "\n", " ", -1)},
	)

	if err != nil {
		return nil, fmt.Errorf("diff command failed: %v", err)
	}

	return strings.Fields(strings.TrimSpace(output)), nil
}

func stepsToTrigger(files []string, watch []WatchConfig) ([]Step, error) {
	steps := []Step{}

	for _, w := range watch {
		for _, p := range w.Paths {
			for _, f := range files {
				match, err := matchPath(p, f)
				if err != nil {
					return nil, err
				}
				if match {
					steps = append(steps, w.Step)
					break
				}
			}
		}
	}

	return dedupSteps(steps), nil
}

// matchPath checks if the file f matches the path p.
func matchPath(p string, f string) (bool, error) {
	// If the path contains a glob, the `doublestar.Match`
	// method is used to determine the match,
	// otherwise `strings.HasPrefix` is used.
	if strings.Contains(p, "*") {
		match, err := doublestar.Match(p, f)
		if err != nil {
			return false, fmt.Errorf("path matching failed: %v", err)
		}
		if match {
			return true, nil
		}
	}
	if strings.HasPrefix(f, p) {
		return true, nil
	}
	return false, nil
}

func dedupSteps(steps []Step) []Step {
	unique := []Step{}
	for _, p := range steps {
		duplicate := false
		for _, t := range unique {
			if reflect.DeepEqual(p, t) {
				duplicate = true
				break
			}
		}

		if !duplicate {
			unique = append(unique, p)
		}
	}

	return unique
}

func generatePipeline(steps []Step, plugin Plugin) (*os.File, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "bmrd-")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary pipeline file: %v", err)
	}

	yamlSteps := make([]yaml.Marshaler, len(steps))

	for i, step := range steps {
		yamlSteps[i] = step
	}

	if plugin.Wait {
		yamlSteps = append(yamlSteps, WaitStep{})
	}

	for _, cmd := range plugin.Hooks {
		yamlSteps = append(yamlSteps, Step{Command: cmd.Command})
	}

	yamlNotify := make([]yaml.Marshaler, len(plugin.Notify))
	for i, n := range plugin.Notify {
		yamlNotify[i] = n
	}

	pipeline := map[string][]yaml.Marshaler{
		"steps": yamlSteps,
	}

	if len(yamlNotify) > 0 {
		pipeline["notify"] = yamlNotify
	}

	data, err := yaml.Marshal(&pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not serialize the pipeline: %v", err)
	}

	// Disable logging in context of go tests.
	if env("TEST_MODE", "") != "true" {
		fmt.Printf("Generated Pipeline:\n%s\n", string(data))
	}

	if err = ioutil.WriteFile(tmp.Name(), data, 0644); err != nil {
		return nil, fmt.Errorf("could not write step to temporary file: %v", err)
	}

	return tmp, nil
}
