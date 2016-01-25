#!/usr/bin/env bats

load ../test_helper

@test "sources: Etcd works" {
  CONTAINER_ID=$(docker run -d --label configo="true" -h consul quay.io/coreos/etcd:v2.0.3 -bind-addr=0.0.0.0:4001)
  for i in {1..5}; do [ "$(docker exec $CONTAINER_ID /etcdctl set myApp/test/property test 2>/dev/null )" = "test" ] && break || sleep 1; done
  
  run_container_with_parameters "--link $CONTAINER_ID:etcd" <<EOC
  export CONFIGO_SOURCE_0='{"type": "etcd", "endpoints": ["http://etcd:4001"], "prefix": "myApp/"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "test"
}

@test "sources: Etcd with KeepPrefix works" {
  CONTAINER_ID=$(docker run -d --label configo="true" -h consul quay.io/coreos/etcd:v2.0.3 -bind-addr=0.0.0.0:4001)
  for i in {1..5}; do [ "$(docker exec $CONTAINER_ID /etcdctl set myApp/test/property test 2>/dev/null )" = "test" ] && break || sleep 1; done
  
  run_container_with_parameters "--link $CONTAINER_ID:etcd" <<EOC
  export CONFIGO_SOURCE_0='{"type": "etcd", "endpoints": ["http://etcd:4001"], "prefix": "myApp/", "keepPrefix": true}'
  configo printenv MYAPP_TEST_PROPERTY
EOC

  assert_success "test"
}