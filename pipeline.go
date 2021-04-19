package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Pipeline is Buildkite pipeline definition
type Pipeline struct {
	Steps []Step
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

	steps := stepsToTrigger(diffOutput, plugin.Watch)

	pipeline, err := generatePipeline(steps, plugin)
	defer os.Remove(pipeline.Name())

	if err != nil {
		log.Error(err)
		return "", []string{}, err
	}

	cmd := "buildkite-agent"
	args := []string{"pipeline", "upload", pipeline.Name()}

	if plugin.Interpolation {
		args = append(args, "--no-interpolation")
	}

	executeCommand("buildkite-agent", args)

	return cmd, args, nil
}

func diff(command string) ([]string, error) {
	log.Infof("Running diff command: %s", command)

	split := strings.Split(command, " ")
	cmd, args := split[0], split[1:]

	output, err := executeCommand(cmd, args)
	if err != nil {
		return nil, fmt.Errorf("diff command failed: %v", err)
	}

	f := func(c rune) bool {
		return c == '\n'
	}

	return strings.FieldsFunc(strings.TrimSpace(output), f), nil
}

func stepsToTrigger(files []string, watch []WatchConfig) []Step {
	steps := []Step{}

	for _, w := range watch {
		for _, p := range w.Paths {
			for _, f := range files {
				if strings.HasPrefix(f, p) {
					steps = append(steps, w.Step)
					break
				}
			}
		}
	}

	return dedupSteps(steps)
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
	pipeline := Pipeline{Steps: steps}

	if err != nil {
		return nil, fmt.Errorf("could not create temporary pipeline file: %v", err)
	}

	data, err := yaml.Marshal(&pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not serialize the pipeline: %v", err)
	}

	if plugin.Wait {
		data = append(data, "- wait"...)
	}

	for _, cmd := range plugin.Hooks {
		data = append(data, "\n- command: "+cmd.Command...)
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
