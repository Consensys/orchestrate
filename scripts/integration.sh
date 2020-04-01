#!/bin/bash

# Exit on error
set -Eeu

mkdir -p build/coverage-integration
go test -tags integration -covermode=count -coverprofile build/coverage-integration/profile.out "$@"

# Ignore generated & testutils files
cat build/coverage-integration/profile.out | grep -Fv -e ".pb.go" -e ".pb.gw.go" -e "/tests" -e "/testutils" -e "/integration-tests"> build/coverage-integration/cover.out

# Generate coverage report in html formart
go tool cover -func=build/coverage-integration/cover.out | grep total:
go tool cover -html=build/coverage-integration/cover.out -o build/coverage-integration/coverage.html

# Remove temporary file
rm build/coverage-integration/profile.out build/coverage-integration/cover.out || true
