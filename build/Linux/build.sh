#!/bin/bash

output_path="$1"
# Write to ./bin if no custom path specified
[ -z "$output_path" ] && output_path="./bin/"
# Ensure the path ends with '/'
[ "${output_path: -1}" != "/" ] && output_path="${output_path}/"
go build -o "${output_path}sdakit"