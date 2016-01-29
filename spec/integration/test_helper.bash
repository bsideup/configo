setup() {
  export PATH=$PATH:$PWD/bin/
}

run_container() {
  run_container_with_parameters "" $1
}

run_container_with_parameters() {
  docker pull buildpack-deps:curl
  run docker run --label configo="true" -i -v "$PWD/bin/configo.linux-amd64":/bin/configo:ro $1 buildpack-deps:curl $2
}

assert_equal() {
  if [ "$1" != "$2" ]; then
    echo "expected: $1"
    echo "actual  : $2"
    return 1
  fi
}

assert_range() {
  if [ $1 -lt $2 ]; then
    echo "expected: $1"
    echo "greater than: $2"
    return 1
  fi
  if [ $1 -gt $3 ]; then
    echo "expected: $1"
    echo "less than: $3"
    return 1
  fi
}

assert_output() {
    assert_equal "$1" "$output"
}

assert_success() {
  if [ "$status" -ne 0 ]; then
    echo "command failed with exit status $status"
    echo "command output: $output"
    return 1
  elif [ "$#" -gt 0 ]; then
    assert_equal "$1" "$output"
  fi
}

assert_failure() {
  if [ "$status" -eq 0 ]; then
    echo "command successed, but should fail"
    echo "command output: $output"
    return 1
  elif [ "$#" -gt 0 ]; then
    assert_equal "$1" "$output"
  fi
}

teardown() {
  docker rm -f $(docker ps -q --filter "configo=true") &>/dev/null | true 
}