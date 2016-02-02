#!/usr/bin/env bats

load ../test_helper

setup() {
  CONTAINER_ID=$(docker run -d --label configo="true" dperson/nginx:latest -g nginx)
  docker exec -i $CONTAINER_ID bash <<EOC
/bin/cat <<EOF >/srv/www/test.json
{
  "some": {
    "nested": {
      "structure": true
    }
  }
}
EOF
EOC
}

@test "sources: HTTP with TLS Cert auth and not configured TLS should fail" {
  docker exec -i $CONTAINER_ID bash <<EOC
sed -i -e 's%ssl on;%ssl on;\n    ssl_client_certificate /etc/nginx/ssl/fullchain.pem;\n    ssl_verify_client on;%g' /etc/nginx/conf.d/default.conf
sleep 1
service nginx reload
EOC

  run_container_with_parameters "--link $CONTAINER_ID:nginx" <<EOC
  export CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "https://nginx/test.json", "insecure": true}'
  configo printenv SOME_NESTED_STRUCTURE
EOC

  assert_failure "400 Bad Request"
}

@test "sources: HTTP with Cert auth works" {
  docker exec -i $CONTAINER_ID bash <<EOC
sed -i -e 's%ssl on;%ssl on;\n    ssl_client_certificate /etc/nginx/ssl/fullchain.pem;\n    ssl_verify_client on;%g' /etc/nginx/conf.d/default.conf
sleep 1
service nginx reload
EOC

  run_container_with_parameters "--link $CONTAINER_ID:nginx --volumes-from=$CONTAINER_ID" <<'EOC'
  export CONFIGO_SOURCE_0=$(cat <<EOF | tr -d "\n"
{
  "type": "http",
  "format": "json",
  "url": "https://nginx/test.json",
  "insecure": true,
  "tls": {
    "cert": "$(while read -r line; do printf "%s\\\n" "$line"; done </etc/nginx/ssl/fullchain.pem)",
    "key": "$(while read -r line; do printf "%s\\\n" "$line"; done </etc/nginx/ssl/privkey.pem)"
  }
}
EOF
)
  configo printenv SOME_NESTED_STRUCTURE
EOC

  assert_success "true"
}

@test "sources: HTTP with insecure source should fail" {
  run_container_with_parameters "--link $CONTAINER_ID:nginx -v /etc/ssl/certs:/etc/ssl/certs" <<EOC
  export CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "https://nginx/test.json"}'
  configo printenv SOME_NESTED_STRUCTURE
EOC

  assert_failure "Get https://nginx/test.json: x509: certificate signed by unknown authority"
}

@test "sources: HTTP with insecure source and insecure:true works" {
  run_container_with_parameters "--link $CONTAINER_ID:nginx" <<EOC
  export CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "https://nginx/test.json", "insecure": true}'
  configo printenv SOME_NESTED_STRUCTURE
EOC

  assert_success "true"
}

@test "sources: HTTP works" {
  run_container_with_parameters "--link $CONTAINER_ID:nginx" <<EOC
  export CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "http://nginx/test.json"}'
  configo printenv SOME_NESTED_STRUCTURE
EOC

  assert_success "true"
}
