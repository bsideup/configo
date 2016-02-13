#!/usr/bin/env bats

load test_helper

@test "encryption: basic" {
  run_container_with_parameters "-e CONFIGO_LOG_LEVEL=ERROR" <<EOC
  export CONFIGO_ENCRYPTION_KEY="a very very very very secret key"
  export TEST_PROPERTY='CONFIGO:{{decrypt "mYii+KwpzEHroZUNuT2jAirM2qmJUr1tdWFFocGEJOQ="}}'
  configo printenv TEST_PROPERTY
EOC
  
  assert_success "123"
}

@test "encryption: fails is encryption key is not set" {
  run_container_with_parameters "-e CONFIGO_LOG_LEVEL=ERROR" <<EOC
  export TEST_PROPERTY='CONFIGO:{{decrypt "mYii+KwpzEHroZUNuT2jAirM2qmJUr1tdWFFocGEJOQ="}}'
  configo printenv TEST_PROPERTY
EOC
  
  assert_failure "template: TEST_PROPERTY:1:2: executing \"TEST_PROPERTY\" at <decrypt \"mYii+KwpzEH...>: error calling decrypt: CONFIGO_ENCRYPTION_KEY should be set in order to use \`decrypt\` function"
}

@test "encryption: wrong Base64" {
  run_container_with_parameters "-e CONFIGO_LOG_LEVEL=ERROR" <<EOC
  export CONFIGO_ENCRYPTION_KEY="a very very very very secret key"
  export TEST_PROPERTY='CONFIGO:{{decrypt "/!qwe123"}}'
  configo printenv TEST_PROPERTY
EOC
  
  assert_failure "template: TEST_PROPERTY:1:2: executing \"TEST_PROPERTY\" at <decrypt \"/!qwe123\">: error calling decrypt: illegal base64 data at input byte 1"
}