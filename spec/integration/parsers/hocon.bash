#!/usr/bin/env bats

load ../test_helper

@test "parsers: HOCON works" {
  run_container <<EOC
  /bin/cat <<EOF >/test.conf
test {
  property = 123
}
EOF

  export CONFIGO_SOURCE_0='type: file, path: /test.conf, format: hocon'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}
