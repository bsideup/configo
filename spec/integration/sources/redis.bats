#!/usr/bin/env bats

load ../test_helper

@test "sources: Redis works" {
  CONTAINER_ID=$(docker run -d --label configo="true" redis:3.0.6-alpine)
  docker run -i --rm --link $CONTAINER_ID:redis redis:3.0.6-alpine redis-cli -h redis hset myAppConfig TEST_PROPERTY 123
  
  run_container_with_parameters "--link $CONTAINER_ID:redis" <<EOC
  export CONFIGO_SOURCE_0='{"type": "redis", "uri": "redis://redis/0", "key": "myAppConfig"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}