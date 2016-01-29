#!/usr/bin/env bats

load ../test_helper

@test "sources: HTTP works" {
  CONTAINER_ID=$(docker run -d --label configo="true" -v /usr/html/ smebberson/alpine-nginx:2.1.1)
  
  run_container_with_parameters "--link $CONTAINER_ID:nginx --volumes-from $CONTAINER_ID" <<EOC
  /bin/cat <<EOF >/usr/html/test.json
{
  "some": {
    "nested": {
      "structure": true
    }
  }
}
EOF

  export CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "http://nginx/test.json"}'
  configo printenv SOME_NESTED_STRUCTURE
EOC

  assert_success "true"
}