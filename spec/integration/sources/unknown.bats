#!/usr/bin/env bats

load ../test_helper

@test "sources: unknown fails" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{"type": "NON_EXISTING_TYPE"}'
  configo printenv TEST_PROPERTY
EOC

  assert_failure "Failed to find source type: NON_EXISTING_TYPE"
}