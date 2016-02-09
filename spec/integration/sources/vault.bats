#!/usr/bin/env bats

load ../test_helper

@test "sources: HTTP with TLS Cert auth and not configured TLS should fail" {  
  CONTAINER_ID=$(
  docker run -d --label configo="true" --entrypoint=/bin/sh -e VAULT_ADDR='http://localhost:8201/' cgswong/vault:0.4.0 -c '
  cat <<EOF > /config.hcl
  listener "tcp" {
   address = "0.0.0.0:8201"
   tls_disable = 1
  }
  disable_mlock = true
EOF

  vault server -dev -config=/config.hcl
')

  docker exec $CONTAINER_ID vault write secret/myapp test_property=true
  
  TOKEN=$(docker exec $CONTAINER_ID vault token-create | grep "token " | tr -s ' ' | cut -f2 | tr -d '\r')

  run_container_with_parameters "--link $CONTAINER_ID:vault" <<EOC
  export CONFIGO_SOURCE_0='{"type": "vault", "address": "http://vault:8201/", "token": "$TOKEN", "path": "secret/myapp"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "true"
}
