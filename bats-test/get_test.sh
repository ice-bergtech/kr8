#!/usr/bin/env bats

if [ -z "$KR8" ]; then
  KR8=kr8
fi

KR8_ARGS="-B data"
CLUSTER=bats

# NOTE: These are expected to be the same as "cluster ..." output, so reuse
# the expected files.  --clusterparams might throw a wrench in this

@test "Check get clusters output" {
  expected=$(<expected/cluster_list)
  run $KR8 $KR8_ARGS get clusters
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "Check get components output" {
  expected=$(<expected/cluster_components)
  run $KR8 $KR8_ARGS get components --cluster "$CLUSTER"
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

# These have a (debug?) output line in the stock version
@test "Check get params for all components" {
  expected=$(<expected/get_params)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER"
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "Check get params for one component with cluster config (-C)" {
  expected=$(<expected/get_params_comp1)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" -c comp1
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "Check get component params for one component only (-P)" {
  expected=$(<expected/get_params_comp2)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" -P comp2
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}
