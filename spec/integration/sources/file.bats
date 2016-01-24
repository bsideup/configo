#!/usr/bin/env bats

load ../test_helper

@test "sources: file works" {
  run_container <<EOC
  /bin/cat <<EOF >/test.yml
test:
  property: 123
EOF

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test.yml", "format": "yaml"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}