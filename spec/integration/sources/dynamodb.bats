#!/usr/bin/env bats

load ../test_helper

@test "sources: DynamoDB works" {
  CONTAINER_ID=$(docker run -d --label configo="true" tutum/dynamodb:latest)
  docker run -i --rm --link $CONTAINER_ID:dynamodb xueshanf/awscli bash <<EOC
export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy

aws dynamodb create-table --endpoint-url http://dynamodb:8000 --region us-west-1 \
  --table-name configs \
  --attribute-definitions AttributeName=key,AttributeType=S \
  --key-schema AttributeName=key,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
    
aws dynamodb put-item --endpoint-url http://dynamodb:8000 --region us-west-1 \
  --table-name configs  \
  --item '{"key":{"S":"myApp" }, "TEST_PROPERTY":{"S": "test"}}'
EOC
  
  run_container_with_parameters "--link $CONTAINER_ID:dynamodb" <<EOC
  export CONFIGO_SOURCE_0='{"type": "dynamodb", "endpoint": "http://dynamodb:8000/", "accessKey":"dummy", "secretKey": "dummy", "table": "configs", "key": "myApp"}'
  configo printenv TEST_PROPERTY
EOC

  assert_success "test"
}