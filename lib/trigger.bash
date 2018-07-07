#!/bin/bash
set -ueo pipefail

function generate_pipeline_yml() {
  for pipeline_index in "${index_of_pipelines_to_trigger[@]}";
    do
      trigger "$pipeline_index"
    done
  add_wait
  add_hooks
}

function trigger() {
  local pipeline=$1
  local trigger
  trigger=$(read_pipeline_config "$pipeline" "TRIGGER")
  echo >&2 "Generating trigger for pipeline: ${trigger}"
  add_trigger "${trigger}"
  add_label "$(read_pipeline_config "$pipeline" "LABEL")"
  add_async "$(read_pipeline_config "$pipeline" "ASYNC")"
  add_branches "$(read_pipeline_config "$pipeline" "BRANCHES")"
  add_build "$pipeline"
}

function add_trigger() {
  local trigger=$1

  if [[ -n $trigger ]];
    then
      pipeline_yml+=("  - trigger: ${trigger}")
    else
      echo "Invalid config. Pipeline trigger is required"
      # exit 1
  fi
}

function add_label() {
  local label=$1

  if [[ -n $label ]]; then
    pipeline_yml+=("    label: ${label}")
  fi
}

function add_async() {
  local async=$1

  if [[ -n $async ]]; then
    pipeline_yml+=("    async: ${async}")
  fi
}

function add_branches() {
  local branches=$1

  if [[ -n $branches ]]; then
    pipeline_yml+=("    branches: ${branches}")
  fi
}

function add_build() {
  local pipeline=$1

  pipeline_yml+=("    build:")
  add_build_commit "$(read_pipeline_config "$pipeline" "BUILD_COMMIT")"
  add_build_message "$(read_pipeline_config "$pipeline" "BUILD_MESSAGE")"
  add_build_branch "$(read_pipeline_config "$pipeline" "BUILD_BRANCH")"
  add_build_env "$pipeline"
}

function add_build_commit() {
  local build_commit=$1
  default_commit=${BUILDKITE_COMMIT:-}

  pipeline_yml+=("      commit: ${build_commit:-$default_commit}")
}

function add_build_message() {
  local build_message=$1
  default_message="${BUILDKITE_MESSAGE:-}"

  pipeline_yml+=("      message: \"${build_message:-$default_message}\"")
}

function add_build_branch() {
  local build_branch=$1
  default_branch=${BUILDKITE_BRANCH:-}

  pipeline_yml+=("      branch: ${build_branch:-$default_branch}")
}

function add_build_env() {
  local pipeline=$1
  local build_env
  build_envs=$(read_pipeline_build_env "$pipeline")

  if [[ -n "$build_envs" ]]; then
    pipeline_yml+=("      env:")
    while IFS=$'\n' read -r build_env ; do
      IFS='=' read -r key value <<< "$build_env"
      if [[ -n "$value" ]]; then
        pipeline_yml+=("        ${key}: ${value}")
      else
        pipeline_yml+=("        ${key}: ${!key}")
      fi
    done <<< "$build_envs"
  fi
}

function add_wait() {
  local wait
  wait=${BUILDKITE_PLUGIN_MONOREPO_DIFF_WAIT:-true}

  if [[ "$wait" = true ]] ; then
    pipeline_yml+=("  - wait")
  fi
}

function add_command() {
  local command=$1

  if [[ -n $command ]];
    then
      pipeline_yml+=("  - command: ${command}")
  fi
}

function add_hooks() {
  while IFS=$'\n' read -r command ; do
    add_command "$command"
  done <<< "$(plugin_read_list HOOKS COMMAND)"
}
