#!/usr/bin/env bats

load ../test_helper

setup() {
  CONTAINER_ID=$(docker run -d --label configo="true" quay.io/coreos/etcd:v2.0.3 -bind-addr=0.0.0.0:4001)
  for i in {1..5}; do [ "$(docker run -i --rm --label configo='true' --link $CONTAINER_ID:etcd --entrypoint=/usr/bin/curl gliderlabs/consul:0.6 -sSL -XPUT http://etcd:4001/v2/keys/myApp/test/property -d value=test 2>/dev/null )" == "{action":"set"* ] && break || sleep 1; done
}

@test "sources: Etcd works" {
  run_container_with_parameters "--link $CONTAINER_ID:etcd" <<EOC
  export CONFIGO_SOURCE_0='type: etcd, endpoints: ["http://etcd:4001"], prefix: "myApp/"'
  configo printenv TEST_PROPERTY
EOC

  assert_success "test"
}

@test "sources: Etcd with KeepPrefix works" {
  run_container_with_parameters "--link $CONTAINER_ID:etcd" <<EOC
  export CONFIGO_SOURCE_0='type: etcd, endpoints: ["http://etcd:4001"], prefix: "myApp/", keepPrefix: true'
  configo printenv MYAPP_TEST_PROPERTY
EOC

  assert_success "test"
}
