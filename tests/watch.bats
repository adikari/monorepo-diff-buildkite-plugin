#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

@test "Notify if the watched path is in the diff output" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "Detected changes in watched path: services/foo"
}

@test "Display the actions triggered by the watched path" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_COMMAND="echo 123"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "Generating trigger for path: services/foo"
  assert_output --partial "Generating command for path: services/bar"
}
