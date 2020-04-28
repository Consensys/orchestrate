#!/bin/bash

# Exit on error
set -Eeu

for p in `find . -name '*.proto' -not -path './vendor/*'`; do
    echo "# $p"

    protoc -I. \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --go_out=plugins=grpc+protoc-gen-grpc-gateway+protoc-gen-swagger,paths=source_relative:. \
      --grpc-gateway_out=logtostderr=true,paths=source_relative:. \
      --swagger_out=logtostderr=true:./public/swagger-specs \
      $p
done
