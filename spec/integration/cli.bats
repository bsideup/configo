#!/usr/bin/env bats

load test_helper

@test "cli: command is required" {
  run_container "configo"
  
  assert_failure "the required argument \`command\` was not provided"
}

@test "cli: command is executed" {
  run_container_with_parameters "-e CONFIGO_LOG_LEVEL=ERROR" <<EOC
  configo echo true
EOC
  
  assert_success "true"
}