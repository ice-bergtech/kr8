#!/bin/bash

set -e

for i in cluster get jsonnet; do
  echo "Generating expected output for '$i' command"
  ./${i}_generate.sh
done
