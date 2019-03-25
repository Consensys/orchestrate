#!/bin/bash

# Exit on error
set -e

echo "mode: count" >> tmp.out
for package in $@; do
  go test -covermode=count -coverprofile profile.out "${package}"
  if [ -f profile.out ]; then
    tail -q -n +2 profile.out >> tmp.out
    rm profile.out
  fi
done

# Generate coverage report in html formart
go tool cover -func=tmp.out
go tool cover -html=tmp.out -o coverage.html

# Remove temporary file
rm tmp.out
