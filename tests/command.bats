#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

@test "Notify when there are no changes" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff3"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"

  run $PWD/hooks/command

  assert_success
  assert_output --partial "No changes detected"
}
