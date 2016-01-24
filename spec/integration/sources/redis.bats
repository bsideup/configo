#!/usr/bin/env bats

load ../test_helper

@test "sources: Redis works" {
  CONTAINER_ID=$(docker run -d --label configo="true" redis)
  docker exec -it $CONTAINER_ID bash -c "redis-cli hset myAppConfig TEST_PROPERTY 123"
  
  run_container_with_parameters "--link $CONTAINER_ID:redis" <<EOC
  export CONFIGO_SOURCE_0='{"type": "redis", "uri": "redis://redis/0", "key": "myAppConfig"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "123"
}