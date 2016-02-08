#!/usr/bin/env bats

load ../test_helper

@test "sources: shell works" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{"type": "shell", "format": "properties", "command": "echo test_property=123 | grep test_"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}