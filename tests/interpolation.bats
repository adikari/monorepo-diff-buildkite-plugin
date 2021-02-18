#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

@test "Runs with interpolation on" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_INTERPOLATION="true"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_COMMAND="echo 123"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload --no-interpolation : echo uploading without interpolation" \
    "pipeline upload : echo uploading with interpolation"

  run $PWD/hooks/command

  unstub buildkite-agent

  assert_success
  assert_output --partial "uploading with interpolation"
}

@test "Runs with interpolation off" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_INTERPOLATION="false"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_COMMAND="echo 123"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload --no-interpolation : echo uploading without interpolation" \
    "pipeline upload : echo uploading with interpolation"

  run $PWD/hooks/command

  unstub buildkite-agent

  assert_success
  assert_output --partial "uploading without interpolation"
}
