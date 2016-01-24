#!/usr/bin/env bats

load ../test_helper

@test "parsers: JSON works" {
  run_container <<EOC
  /bin/cat <<EOF >/test.json
{
  "test": {
    "property": 123
  }
}
EOF

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test.json", "format": "json"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}