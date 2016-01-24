#!/usr/bin/env bats

load ../test_helper

@test "parsers: Properties works" {
  run_container <<EOC
  /bin/cat <<EOF >/test.properties
test_property=123
EOF

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test.properties", "format": "properties"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}