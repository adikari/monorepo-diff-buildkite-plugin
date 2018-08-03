#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

@test "Generate steps" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "steps:"
}

@test "Generates trigger" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_LABEL="Bar service deployment"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "  - trigger: slug-for-foo"
  assert_output --partial "  - trigger: slug-for-bar"
}

@test "Uses user defined label if it exist" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_LABEL="Bar service deployment"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "    label: Bar service deployment"
}

@test "Adds async if supplied" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_ASYNC="true"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "    async: true"
}

@test "Omits async if not supplied" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  refute_output --partial "    async: true"
}

@test "Adds branches if supplied" {
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BRANCHES="!master"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_ASYNC="true"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "    branches: !master"
}

@test "Adds env if supplied" {
  export AWS_REGION="us-east-1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BRANCHES="!master"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_ASYNC="true"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_0="AWS_REGION"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_1="NODE_ENV=production"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "    env:"
  assert_output --partial "      AWS_REGION: us-east-1"
  assert_output --partial "      NODE_ENV: production"
}

@test "Adds wait at the end of triggers by default" {
  export AWS_REGION="us-east-1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BRANCHES="!master"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_ASYNC="true"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_0="AWS_REGION"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_1="NODE_ENV=production"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "  - wait"
}

@test "Does not wait on triggers if explicitly set to false" {
  export AWS_REGION="us-east-1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BRANCHES="!master"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_ASYNC="true"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_0="AWS_REGION"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_1="NODE_ENV=production"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WAIT=false
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  refute_output --partial "  - wait"
}

@test "Adds hook commands" {
  export AWS_REGION="us-east-1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BRANCHES="!master"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_ASYNC="true"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_0="AWS_REGION"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_BUILD_ENV_1="NODE_ENV=production"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_PATH="services/path-not-in-diff"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_1_CONFIG_TRIGGER="slug-for-path-not-in-diff"

  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_PATH="services/bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_2_CONFIG_TRIGGER="slug-for-bar"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_HOOKS_0_COMMAND="echo test"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_HOOKS_1_COMMAND="aws s3 ls"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial "  - command: echo test"
  assert_output --partial "  - command: aws s3 ls"
}

@test "Preserves quotes in commit messages" {
  export BUILDKITE_MESSAGE="Hello world \"stuff\" \n multiline"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF="cat $PWD/tests/mocks/diff1"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_PATH="services/foo"
  export BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_0_CONFIG_TRIGGER="slug-for-foo"
  export DEBUG=true

  stub buildkite-agent \
    "pipeline upload : echo uploading"

  run $PWD/hooks/command

  unstub buildkite-agent
  assert_success
  assert_output --partial '    message: "Hello world \"stuff\" \n multiline"'
}
