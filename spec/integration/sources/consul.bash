#!/usr/bin/env bats

load ../test_helper

@test "sources: Consul works" {
  CONTAINER_ID=$(docker run -d --label configo="true" gliderlabs/consul:0.6 agent -dev -client=0.0.0.0)
  for i in {1..5}; do [ "$(docker run --label configo="true" -i --rm --link $CONTAINER_ID:consul --entrypoint=/usr/bin/curl gliderlabs/consul:0.6 -sSL -X PUT -d 'test' http://consul:8500/v1/kv/myAppConfig/TEST_PROPERTY)" = "true" ] && break || sleep 1; done
  
  run_container_with_parameters "--link $CONTAINER_ID:consul" <<EOC
  export CONFIGO_SOURCE_0='type: consul, address: "consul:8500", scheme: http, prefix: myAppConfig'
  configo printenv TEST_PROPERTY
EOC

  assert_success "test"
}
