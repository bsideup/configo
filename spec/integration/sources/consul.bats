#!/usr/bin/env bats

load ../test_helper

@test "sources: Consul works" {
  CONTAINER_ID=$(docker run -d --label configo="true" -h consul progrium/consul -server -bootstrap)
  until [ "$(docker exec $CONTAINER_ID bash -c "curl -sSL -X PUT -d 'test' http://localhost:8500/v1/kv/myAppConfig/TEST_PROPERTY")" = "true" ]; do
    sleep 1;
  done
  
  run_container_with_parameters "--link $CONTAINER_ID:consul" <<EOC
  export CONFIGO_SOURCE_0='{"type": "consul", "address": "consul:8500", "scheme": "http", "prefix": "myAppConfig"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "test"
}