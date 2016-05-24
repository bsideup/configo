#!/usr/bin/env bats

load ../test_helper

@test "sources: unknown fails" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{"type": "NON_EXISTING_TYPE"}'
  configo printenv TEST_PROPERTY
EOC

  assert_failure "Failed to parse source: Failed to find source type: NON_EXISTING_TYPE"
}

@test "sources: should fail on unknown field" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{"type": "http", "fomat": "json"}'
  configo env
EOC

  assert_failure "Failed to parse source: unknown configuration keys: [fomat]"
}