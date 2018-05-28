#!/bin/bash
set -euo pipefail

function run_diff() {
  local diff_command
  local default_diff_command
  local diff_output

  default_diff_command="git diff --name-only HEAD~1"
  diff_command="${BUILDKITE_PLUGIN_MONOREPO_DIFF_DIFF:-$default_diff_command}"
  diff_output=$($diff_command)

  echo >&2 "Running diff command: [$diff_command]" ;
  echo >&2 "Output from diff:" ;
  echo >&2 "$diff_output";

  echo "$diff_output"
}