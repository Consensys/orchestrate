#!/bin/bash

# Exit on error
set -Eeu

echo "mode: count" > tmp.out
for package in $@; do
  go test -covermode=count -coverprofile profile.out "${package}"
  if [ -f profile.out ]; then
    tail -q -n +2 profile.out >> tmp.out
    rm profile.out
  fi
done

# Ignore generated files
cat tmp.out | grep -v ".pb.go" --exclude-dir=examples --exclude-dir=e2e > cover.out

# Generate coverage report in html formart
go tool cover -func=cover.out
go tool cover -html=cover.out -o coverage.html

# Remove temporary file
rm tmp.out cover.out
