#!/bin/bash

# Reads either a value or a list from plugin config
function plugin_read_list() {
  local prefix="BUILDKITE_PLUGIN_MONOREPO_DIFF_$1"
  local suffix="$2"
  local parameter="${prefix}_0_${suffix}"

  if [[ -n "${!parameter:-}" ]]; then
    local i=0
    local parameter="${prefix}_${i}_${suffix}"
    while [[ -n "${!parameter:-}" ]]; do
      echo "${!parameter}"
      i=$((i+1))
      parameter="${prefix}_${i}_${suffix}"
    done
  elif [[ -n "${!prefix:-}" ]]; then
    echo "${!prefix}"
  fi
}

function read_pipeline_config() {
  local pipeline_index=$1
  local config_key=$2
  local parameter="BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_${pipeline_index}_CONFIG_${config_key}"

  echo "${!parameter:-}"
}

function read_pipeline_build_env() {
  local pipeline_index=$1
  local prefix="BUILDKITE_PLUGIN_MONOREPO_DIFF_WATCH_${pipeline_index}_CONFIG_BUILD_ENV"
  local parameter="${prefix}_0"

  if [[ -n "${!parameter:-}" ]]; then
    local i=0
    local parameter="${prefix}_${i}"
    while [[ -n "${!parameter:-}" ]]; do
      echo "${!parameter}"
      i=$((i+1))
      parameter="${prefix}_${i}"
    done
  fi
}

function in_array() {
  local e
  for e in "${@:2}"; do [[ "$e" == "$1" ]] && return 0; done
  return 1
}