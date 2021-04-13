#!/usr/bin/env bats

load '/usr/local/lib/bats/load.bash'

@test "Notify when there are no changes" {
  run $PWD/hooks/command

  assert_output --partial "No changes detected"
}
