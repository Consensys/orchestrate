FROM golang:1.16.9

ARG PROTOC_TAG
ARG PROTOC_GEN_GO_TAG

RUN apt-get update && apt-get install -y unzip

WORKDIR /protoc

RUN curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_TAG}/protoc-${PROTOC_TAG}-linux-x86_64.zip -o protoc.zip \
    && unzip protoc.zip \
    && mv ./bin/protoc /usr/local/bin/protoc \
    && mv ./include/google /usr/local/include/google \
    && rm -rf /protoc

RUN GO111MODULE=on go get \
    google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_TAG}

WORKDIR /src

USER 1000
