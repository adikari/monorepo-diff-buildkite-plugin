package main

import (
	"errors"
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
type PipelineGenerator func(steps []Step) (*os.File, error)

func uploadPipeline(plugin Plugin, generatePipeline PipelineGenerator) (string, []string, error) {
	diffOutput := diff(plugin.Diff)

	if len(diffOutput) < 1 {
		log.Info("No changes detected. Skipping pipeline upload.")
		return "", []string{}, nil
	}

	steps := stepsToTrigger(diffOutput, plugin.Watch)

	pipeline, err := generatePipeline(steps)
	defer os.Remove(pipeline.Name())

	if err != nil {
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

func diff(command string) []string {
	log.Infof("Running diff command: %s", command)

	split := strings.Split(command, " ")
	cmd, args := split[0], split[1:]

	output, err := executeCommand(cmd, args)

	if err != nil {
		log.Fatalf("%s: %s", err, command)
	}

	log.Debug("Output from diff:\n" + output)

	return strings.Split(strings.TrimSpace(output), "\n")
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

func generatePipeline(steps []Step) (*os.File, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "bmrd-")
	pipeline := Pipeline{Steps: steps}

	if err != nil {
		log.Debug(err)
		return nil, errors.New("Could not create temporary pipeline file")
	}

	data, err := yaml.Marshal(&pipeline)

	log.Debugf("Pipeline to upload:\n%s", string(data))

	if err = ioutil.WriteFile(tmp.Name(), data, 0644); err != nil {
		log.Debug(err)
		return nil, errors.New("Could not write step to temporary file")
	}

	return tmp, nil
}
