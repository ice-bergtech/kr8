#!/bin/bash

for i in cluster get init jsonnet render yaml; do
  echo "Testing '$i' command"
  ./bats/bin/bats ${i}_test.sh
  echo
done
