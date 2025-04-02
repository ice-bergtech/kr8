#!/usr/bin/env bats

if [ -z "$KR8" ]; then
  KR8=kr8
fi

KR8_ARGS="-B data"
CLUSTER=bats

@test "01 Check cluster list output" {
  expected=$(<expected/cluster_list)
  run $KR8 $KR8_ARGS get clusters
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "02 Check cluster components output" {
  expected=$(<expected/cluster_components)
  run $KR8 $KR8_ARGS get components --cluster "$CLUSTER"
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

## The params tests also effectively test param hierarchy at all levels with "comp2"
## FIXME: Params above the cluster/<x>/params.jsonnet hierarchy bleed into the
##        params no matter which component is requested.  Is this correct?
##        "comp1" displays this behavior in a test in case it changes

@test "03 Check cluster params for all components" {
  expected=$(<expected/cluster_params)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER"
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "04 Check cluster params for one component with cluster config (-c)" {
  expected=$(<expected/cluster_params_comp1)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" -c comp1
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "05 Check cluster params for one component only (-P)" {
  expected=$(<expected/cluster_params_comp2)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" -P comp2
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "06 Check cluster params with file override" {
  expected=$(<expected/cluster_params_file)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" --clusterparams data/misc/cluster_params.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

# Check behavior on a component that doesn't exist
@test "07 Check cluster params with unset component (-P)" {
  expected=""
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" -P no_component
  [ "$status" -eq 1 ]
}

# Check behavior on a component that doesn't exist
@test "08 Check cluster params with unset component (-C)" {
  expected=$(<expected/cluster_params_no_comp)
  run $KR8 $KR8_ARGS get params -C "$CLUSTER" -c no_component
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}
