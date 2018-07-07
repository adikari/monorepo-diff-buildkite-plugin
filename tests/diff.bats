#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

@test "Run the specified inline diff command" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "services/foo/serverless.yml"
  assert_output --partial "services/bar/config.yml"
  assert_output --partial "ops/bar/config.yml"
  assert_output --partial "README.md"
}

@test "Run the specified shell script" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="$PWD/tests/mocks/diff2.sh"

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "diff 2"
}

@test "Run default diff command when shell script is not specified" {
  stub git \
    "diff --name-only HEAD~1 : echo services/git/head-minus-1.yml"
  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  unstub git
  assert_success
  assert_output --partial "services/git/head-minus-1.yml"
}

@test "Exits on shell script error" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="$PWD/tests/mocks/thisdoesnotexist.sh"

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  [ "$status" -eq 1 ]
  assert_output --partial "Failed to run diff command"
}