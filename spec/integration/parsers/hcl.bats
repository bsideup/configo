#!/usr/bin/env bats

load ../test_helper

@test "parsers: HCL works" {
  run_container <<EOC
  /bin/cat <<EOF >/test.hcl
test_property = 123
EOF

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test.hcl", "format": "hcl"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}