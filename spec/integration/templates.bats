#!/usr/bin/env bats

load test_helper

@test "templates: functions" {
  run_container_with_parameters "-e CONFIGO_LOG_LEVEL=ERROR" <<EOC
  export SOME_PROPERTY=123
  export TEST_PROPERTY="CONFIGO:some property value is: {{or .NON_EXISTING_PROPERTY .ANOTHER_NON_EXISTING_PROPERTY .SOME_PROPERTY `default`}}"
  configo printenv TEST_PROPERTY
EOC
  
  assert_success "some property value is: 123"
}

@test "templates: fromJSON function" {
  run_container_with_parameters "-e CONFIGO_LOG_LEVEL=ERROR" <<EOC
  export SOME_JSON='{"test": {"inner": 123}}'
  export TEST_PROPERTY="CONFIGO:inner value is: {{(fromJSON .SOME_JSON).test.inner}}"
  configo printenv TEST_PROPERTY
EOC

  assert_success "inner value is: 123"
}