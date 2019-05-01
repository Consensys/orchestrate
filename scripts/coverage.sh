#!/bin/bash

# Exit on error
set -e

# Make sure path is correct
if [ ! -f "scripts/coverage.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create temporary directory
tmpdir=`mktemp -d`

for package in $@; do
  tmpfile=`mktemp -p ${tmpdir} --suffix .cov.tmp`
  go test -covermode=count -coverprofile "${tmpfile}" "${package}"
done

tmpfile=`mktemp -p ${tmpdir} --suffix .cov`
echo "mode: count" > "${tmpfile}"
tail -q -n +2 "${tmpdir}/"*.cov.tmp >> "${tmpfile}"

go tool cover -func="${tmpfile}"

# Generate coverage report in html formart
go tool cover -html="${tmpfile}" -o coverage.html

# Remove tempdir
rm -R "${tmpdir}"
