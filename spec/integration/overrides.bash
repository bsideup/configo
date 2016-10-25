#!/usr/bin/env bats

load test_helper

@test "overrides: multiple sources works" {
  run_container <<'EOC'
  /bin/cat <<EOF >/test1.yml
test_property: 123
another_property: true
EOF

  /bin/cat <<EOF >/test2.yml
another_property: false
third_property: Hi
EOF

  echo "test_property: overrided" >/test3.yml

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test1.yml", "format": "yaml"}'
  export CONFIGO_SOURCE_1='{"type": "file", "path": "/test2.yml", "format": "yaml"}'
  export CONFIGO_SOURCE_2='{"type": "file", "path": "/test3.yml", "format": "yaml"}'
  configo 'echo $TEST_PROPERTY $ANOTHER_PROPERTY $THIRD_PROPERTY'
EOC

  assert_success "overrided false Hi"
}

@test "overrides: ignores nulls" {
  run_container <<'EOC'
  /bin/cat <<EOF >/test1.yml
test_property: 123
another_property: true
EOF

  /bin/cat <<EOF >/test2.yml
another_property: null
EOF

  export CONFIGO_SOURCE_0='{"type": "file", "path": "/test1.yml", "format": "yaml"}'
  export CONFIGO_SOURCE_1='{"type": "file", "path": "/test2.yml", "format": "yaml"}'
  configo 'echo $TEST_PROPERTY $ANOTHER_PROPERTY'
EOC

  assert_success "123 true"
}
