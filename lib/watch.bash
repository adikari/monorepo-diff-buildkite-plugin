#!/bin/bash
set -ueo pipefail

function has_changed() {
  local diff_output=$1
  local watched_path=$2
  for path in $watched_path; do
    if grep -q "^$path" <<< "$diff_output"; then
      return 0
    fi
  done
  return 1
}

function populate_pipelines_to_trigger() {
  local diff_output=$1

  watch_index=0
  while IFS=$'\n' read -r watched_path ; do
    echo >&2 "Comparing watch path: ${watched_path}"
    if has_changed "$diff_output" "$watched_path" ; then
      echo >&2 "Detected changes in watched path: $watched_path"
      pipelines_to_trigger+=("${watch_index} ${watched_path}")
    fi
    watch_index=$((watch_index+1))
  done <<< "$(plugin_read_list WATCH PATH)"
}
