#!/usr/bin/env bats

if [ -z "$KR8" ]; then
  KR8=kr8
fi

@test "01 Check init success" {
  rm -rf ./init-test
  run $KR8 init repo ./init-test
  [ "$status" -eq 0 ]
  [ -d "init-test/clusters" ]
  [ -d "init-test/components" ]
  [ -d "init-test/lib" ]
  rm -rf ./init-test
}

@test "02 Check init success - fetch libs" {
  rm -rf ./init-test
  run $KR8 init repo -f true ./init-test
  [ "$status" -eq 0 ]
  [ -d "init-test/clusters" ]
  [ -d "init-test/components" ]
  [ -d "init-test/lib" ]
  [ -d "init-test/lib/klib" ]
  rm -rf ./init-test
}

@test "03 Check init cluster - named cluster" {
  expected=$(<expected/init_cluster_temp-cluster)
  mkdir -p ./init-test-cluster
  run $KR8 init cluster -B ./init-test-cluster -o temp-cluster
  [ "$status" -eq 0 ]
  diff <(cat "init-test-cluster/clusters/temp-cluster/cluster.jsonnet") <(echo "$expected")
  [ -d "init-test-cluster/clusters/" ]
  rm -rf ./init-test-cluster
}

@test "04 Check init component - named component" {
  expected=$(<expected/init_component_temp-component)
  mkdir -p ./init-test-component
  run $KR8 init component -B ./init-test-component -o temp-component
  [ "$status" -eq 0 ]
  [ -d "init-test-component/components/temp-component" ]
  diff <(cat "init-test-component/components/temp-component/params.jsonnet") <(echo "$expected")
  rm -rf ./init-test-component
}

@test "05 Check init failure - existing directory" {
  skip "skip testing, code issue"
  mkdir -p ./init-test2
  run $KR8 init repo ./init-test2
  [ "$status" -eq 1 ]
  rm -rf ./init-test2
}
