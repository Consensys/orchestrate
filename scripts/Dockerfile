FROM golang:1.12

ARG PROTOC_TAG
ARG PROTOC_GEN_GO_TAG
ARG GRPC_GATEWAY_TAG

RUN apt-get update && apt-get install unzip

WORKDIR /protoc

RUN curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_TAG}/protoc-${PROTOC_TAG}-linux-x86_64.zip -o protoc.zip

RUN unzip protoc.zip

RUN mv ./bin/protoc  /usr/local/bin/protoc && mv ./include/google /usr/local/include/google

WORKDIR /

RUN rm -rf protoc

RUN GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@${PROTOC_GEN_GO_TAG} \
       github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@${GRPC_GATEWAY_TAG} \
       github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@${GRPC_GATEWAY_TAG}

WORKDIR /src