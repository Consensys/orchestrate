#!/bin/bash

 for p in `find . -name *.proto`; do
    echo "# $p"
    protoc -I. --go_out=paths=source_relative:. $p
 done