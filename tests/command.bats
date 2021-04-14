#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

setup() {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_BUILDKITE_PLUGIN_TEST_MODE="true"

  stub buildkite-agent pipeline upload
}

@test "Notify when the plugin cannot be parsed" {

  run $PWD/hooks/command

  assert_failure
  assert_output --partial "Failed to parse plugin configuration"
}

@test "Checks uploaded projects" {
  DIFF_CMD="echo foo-service/"
  LOG_LEVEL="debug"

  export BUILDKITE_PLUGINS="[{\"github.com/chronotc/monorepo-diff-buildkite-plugin\":{\"diff\":\"$DIFF_CMD\",\"log_level\": \"$LOG_LEVEL\",\"wait\":true,\"watch\":[{\"path\":\"foo-service/\",\"config\":{\"trigger\":\"foo-service\"}},{\"path\":\"hello-service/\",\"config\":{\"trigger\":\"this-pipeline-does-not-exists\"}}]}}]"

  run $PWD/hooks/command

  assert_success

  assert_output --partial "Output from diff: \nfoo-service/"

  assert_output --partial << EOM
steps:
  - trigger: foo-service
    build:
      commit: 2569aeab70f9c3f1c04cd436aca7ac501b6ee71b
      message: "fix: temp file not correctly deleted"
      branch: go-rewrite
  - wait
EOM
}
