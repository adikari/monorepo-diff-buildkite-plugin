#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

# https://buildkite.com/kuda/monorepo-diff-buildkite-plugin/builds/110#2b88547e-eb5a-4f99-a83b-affdbcddb303

setup() {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_BUILDKITE_PLUGIN_TEST_MODE="true"

  stub buildkite-agent pipeline upload
}

@test "Notify when the plugin cannot be parsed" {

  run $PWD/hooks/command

  assert_failure
  assert_output --partial "Failed to parse plugin configuration"
}

@test "Check pipeline is generated" {
  DIFF_CMD="echo foo-service/"
  LOG_LEVEL="debug"

  export BUILDKITE_PLUGINS="[{\"github.com/chronotc/monorepo-diff-buildkite-plugin\":{\"diff\":\"$DIFF_CMD\",\"log_level\": \"$LOG_LEVEL\",\"wait\":true,\"watch\":[{\"path\":\"foo-service/\",\"config\":{\"trigger\":\"foo-service\"}},{\"path\":\"hello-service/\",\"config\":{\"trigger\":\"hello-service\"}}]}}]"

  run $PWD/hooks/command

  assert_success

  assert_output --partial "Output from diff: \nfoo-service/"

  assert_output --partial << EOM
steps:
  - trigger: foo-service
  - wait
EOM
}
