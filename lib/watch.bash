#!/bin/bash
set -ueo pipefail

function has_changed() {
  local diff_output=$1
  local watched_path=$2

  if echo "$diff_output" | grep -q "^$watched_path" ; then
    return 0
  else
    return 1
  fi
}

function get_index_of_pipelines_to_trigger() {
  local diff_output=$1

  watch_index=0
  while IFS=$'\n' read -r watched_path ; do
    echo >&2 "Comparing watch path: ${watched_path}"
    if has_changed "$diff_output" "$watched_path" ; then
      echo >&2 "Detected changes in watched path: $watched_path"
      index_of_pipelines_to_trigger+=("$watch_index")
    fi
    watch_index=$((watch_index+1))
  done <<< "$(plugin_read_list WATCH PATH)"
}
