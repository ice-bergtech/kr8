#!/bin/bash

# to setup, run:
# git submodule update --remote --init

for i in cluster get init jsonnet render; do
  echo "Testing '$i' command"
  ./bats/bin/bats ${i}_test.sh
  echo
done
