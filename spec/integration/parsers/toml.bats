#!/usr/bin/env bats

load ../test_helper

@test "parsers: TOML works" {
  run_container <<EOC
  /bin/cat <<EOF >/test.toml
[test]
property = 123
EOF

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test.toml", "format": "toml"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}