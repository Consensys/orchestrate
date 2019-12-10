#!/bin/bash

# Exit on error
set -Eeu

for p in `find . -name *.proto`; do
    echo "# $p"

    protoc -I. \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --go_out=plugins=grpc+protoc-gen-grpc-gateway+protoc-gen-swagger,paths=source_relative:. \
      $p

    protoc -I. \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --grpc-gateway_out=logtostderr=true,paths=source_relative:. \
      $p

    mkdir -p ./public/swagger-specs
    protoc -I. \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --swagger_out=logtostderr=true:./public/swagger-specs \
      $p
done
