#!/usr/bin/env bats

load test_helper

@test "cli: command is required" {
  run_container "configo"
  
  assert_failure
  assert_output "the required argument \`command\` was not provided"
}

@test "cli: command is executed" {
  run_container <<EOC
  configo echo true
EOC
  
  assert_success "true"
}