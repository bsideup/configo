#!/usr/bin/env bats

load ../test_helper

@test "sources: composite works" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='
    type: composite
    sources: [
        { type: shell, format: properties, command: "echo MY_COOL_PROP_1=123" }
        { type: shell, format: properties, command: "echo MY_COOL_PROP_1=override" } 
        { type: shell, format: properties, command: "echo MY_COOL_PROP_2=456" }
    ]
'
  configo printenv MY_COOL_PROP_1
  configo printenv MY_COOL_PROP_2
EOC

  assert_success "override
456"
}

@test "sources: composite nested sources works" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='
    type: composite
    sources: [
        {
            type: composite
            sources: [
                { type: shell, format: properties, command: "echo MY_COOL_PROP_1=123" } 
            ]
        }
        { type: shell, format: properties, command: "echo MY_COOL_PROP_2=456" }
    ]
}'
  configo printenv MY_COOL_PROP_1
  configo printenv MY_COOL_PROP_2
EOC

  assert_success "123
456"
}
