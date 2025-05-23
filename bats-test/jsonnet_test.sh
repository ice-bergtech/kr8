#!/usr/bin/env bats

# FIXME: Add --prune tests

if [ -z "$KR8" ]; then
  KR8=kr8
fi

KR8_ARGS="-B data"
CLUSTER=bats

@test "01 Check jsonnet json parsing" {
  expected=$(<expected/jsonnet_basic_json)
  run $KR8 $KR8_ARGS jsonnet render -C $CLUSTER data/misc/basic.json
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "02 Check jsonnet basic jsonnet parsing" {
  expected=$(<expected/jsonnet_basic_jsonnet)
  run $KR8 $KR8_ARGS jsonnet render -C $CLUSTER data/misc/basic.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "03 Check jsonnet component parsing (default: json)" {
  expected=$(<expected/jsonnet_comp1_json)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 data/components/comp1/comp1.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

# this is a bug where we stacktrace if --component isn't set
# FIXME: could be better
@test "04 Check jsonnet parsing without component - FAIL" {
  expected=$(<expected/jsonnet_comp1_json)
  run $KR8 $KR8_ARGS jsonnet render -C bats data/components/comp1/comp1.jsonnet
  [ "$status" -eq 1 ]
}

# Explicit formats
@test "05 Check jsonnet component parsing (format: json)" {
  expected=$(<expected/jsonnet_comp1_json)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 -F json data/components/comp1/comp1.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "06 Check jsonnet component parsing (format: yaml)" {
  expected=$(<expected/jsonnet_comp1_yaml)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 -F yaml data/components/comp1/comp1.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

# stream format with one object is a stacktrace
# FIXME: could be better
@test "07 Check jsonnet component parsing (format: stream) - FAIL" {
  expected=$(<expected/jsonnet_comp1_json)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 -F stream data/components/comp1/comp1.jsonnet
  [ "$status" -eq 1 ]
}

# List of objects
@test "08 Check jsonnet list component parsing (format: json)" {
  expected=$(<expected/jsonnet_comp1_list_json)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 -F json data/components/comp1/comp1_list.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "09 Check jsonnet list component parsing (format: yaml)" {
  expected=$(<expected/jsonnet_comp1_list_yaml)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 -F yaml data/components/comp1/comp1_list.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

@test "10 Check jsonnet list component parsing (format: stream)" {
  expected=$(<expected/jsonnet_comp1_list_stream)
  run $KR8 $KR8_ARGS jsonnet render -C bats -c comp1 -F stream data/components/comp1/comp1_list.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

# Test with --clusterparams
@test "11 Check jsonnet parsing with --clusterparams" {
  expected=$(<expected/jsonnet_comp2_with_file_stream)
  run $KR8 $KR8_ARGS jsonnet render -c comp2 -F yaml data/components/comp2/comp2.jsonnet \
    --clusterparams data/misc/cluster_params.jsonnet
  [ "$status" -eq 0 ]
  diff <(echo "$output") <(echo "$expected")
}

# FIXME: stacktrace if we call a component that doesn't exist in the --clusterparams file
#        even if that component exists and has its own params
#        Only the clusterprams file gets used, even blanking other params
@test "12 Check jsonnet parsing with --clusterparams FAIL" {
  #expected=$(<expected/jsonnet_comp2_with_file_stream)
  run $KR8 $KR8_ARGS jsonnet render -c comp1 -F stream data/components/comp2/comp1_list.jsonnet \
    --clusterparams data/misc/cluster_params.jsonnet
  [ "$status" -eq 1 ]
  #diff <(echo "$output") <(echo "$expected")
}
