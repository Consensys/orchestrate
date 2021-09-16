#!/bin/bash

# Exit on error
set -Eeu

for p in `find . -name '*.proto' -not -path './vendor/*'`; do
    echo "# $p"
    
    protoc -I. --go_out=paths=source_relative:. $p
done
