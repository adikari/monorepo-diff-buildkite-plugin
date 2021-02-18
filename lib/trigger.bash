#!/bin/bash
set -ueo pipefail

function generate_pipeline_yml() {
  for pipeline in "${pipelines_to_trigger[@]}";
    do
      set -- $pipeline
      # Word split on purpose using spaces

      local pipeline_index=$1
      local pipeline_path=$2
      add_action "$pipeline_index" "$pipeline_path"
    done
  add_wait
  add_hooks
}

function add_action() {
  local pipeline_index=$1
  local pipeline_path=$2

  local action_trigger
  local command

  action_trigger=$(read_pipeline_config "$pipeline_index" "TRIGGER")
  action_command=$(read_pipeline_config "$pipeline_index" "COMMAND")

  if [[ -n $action_trigger ]]
  then
    add_action_trigger "$pipeline_index" "$pipeline_path"
  elif [[ -n $action_command ]]
  then
    add_action_command "$pipeline_index" "$pipeline_path"
  else
    echo "Invalid config. Pipeline trigger or command is required"
  fi
}

function add_action_command() {
  local pipeline_index=$1
  local pipeline_path=$2

  echo >&2 "Generating command for path: $pipeline_path"

  local command
  command=$(read_pipeline_config "$pipeline_index" "COMMAND")

  pipeline_yml+=("  - command: ${command}")

  add_label "$(read_pipeline_config "$pipeline_index" "LABEL")"
  add_agents "$pipeline_index"
  add_artifacts "$pipeline_index"
}

function add_action_trigger() {
  local pipeline_index=$1
  local pipeline_path=$2

  echo >&2 "Generating trigger for path: $pipeline_path"

  local trigger
  trigger=$(read_pipeline_config "$pipeline_index" "TRIGGER")

  pipeline_yml+=("  - trigger: ${trigger}")

  add_label "$(read_pipeline_config "$pipeline_index" "LABEL")"
  add_async "$(read_pipeline_config "$pipeline_index" "ASYNC")"
  add_branches "$(read_pipeline_config "$pipeline_index" "BRANCHES")"
  add_build "$pipeline_index"
}

function add_label() {
  local label=$1

  if [[ -n $label ]]; then
    pipeline_yml+=("    label: \"${label}\"")
  fi
}

function add_agents() {
  local pipeline=$1

  pipeline_yml+=("    agents:")
  add_agents_queue "$(read_pipeline_config "$pipeline" "AGENTS_QUEUE")"
}

function add_artifacts() {
  local pipeline=$1

  local prefix="BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_${pipeline}_CONFIG_ARTIFACTS"
  local parameter="${prefix}_0"

  pipeline_yml+=("    artifacts:")
  if [[ -n "${!parameter:-}" ]]; then
    local i=0
    local parameter="${prefix}_${i}"
    while [[ -n "${!parameter:-}" ]]; do
      pipeline_yml+=("        - ${!parameter}")
      i=$((i+1))
      parameter="${prefix}_${i}"
    done
  fi
}

function add_agents_queue() {
  local queue=$1

  if [[ -n $queue ]]; then
    pipeline_yml+=("      queue: ${queue}")
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
  sanitized_build_message=$(sanitize_string "${build_message:-$default_message}")

  pipeline_yml+=("      message: \"$sanitized_build_message\"")
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

function sanitize_string() {
  local string=$1
  escaped_quotes="${string//\"/\\\"}"
  echo "$escaped_quotes"
}
