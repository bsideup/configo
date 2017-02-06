#!/usr/bin/env bats

load ../test_helper

@test "sources: composite works" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{
    "type": "composite",
    "sources": [
        {
            "type": "shell",
            "format": "properties",
            "command": "echo MY_COOL_PROP_1=123"
        }, 
        {
            "type": "shell",
            "format": "properties",
            "command": "echo MY_COOL_PROP_1=override"
        }, 
        {
            "type": "shell",
            "format": "properties",
            "command": "echo MY_COOL_PROP_2=456"
        }
    ]
}'
  configo printenv MY_COOL_PROP_1
  configo printenv MY_COOL_PROP_2
EOC

  assert_success "override
456"
}

@test "sources: composite nested sources works" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{
    "type": "composite",
    "sources": [
        {
            "type": "composite",
            "sources": [
                {
                    "type": "shell",
                    "format": "properties",
                    "command": "echo MY_COOL_PROP_1=123"
                }
            ]
        },
        {
            "type": "shell",
            "format": "properties",
            "command": "echo MY_COOL_PROP_2=456"
        }
    ]
}'
  configo printenv MY_COOL_PROP_1
  configo printenv MY_COOL_PROP_2
EOC

  assert_success "123
456"
}

@test "sources: composite with uppercasing keys" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{
    "type": "composite",
    "sources": [
        {
            "type": "shell",
            "format": "properties",
            "command": "echo MY_COOL_PROP_1=123"
        },
        {
            "type": "shell",
            "format": "properties",
            "command": "echo my_cool_prop_1=override"
        }
    ]
}'
  configo printenv MY_COOL_PROP_1
EOC

  assert_success "override"
}

@test "sources: composite without uppercasing keys" {
  run_container <<EOC
  export CONFIGO_SOURCE_0='{
    "type": "composite",
    "sources": [
        {
            "type": "shell",
            "format": "properties",
            "command": "echo MY_COOL_PROP_1=123"
        },
        {
            "type": "shell",
            "format": "properties",
            "command": "echo my_cool_prop_1=override"
        }
    ]
}'
  export CONFIGO_UPPERCASE_KEYS=0
  configo printenv MY_COOL_PROP_1
  configo printenv my_cool_prop_1
EOC

  assert_success "123
override"
}
