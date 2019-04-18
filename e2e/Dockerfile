FROM ubuntu:16.04

RUN apt-get update -y && \
    apt-get install -y software-properties-common jq curl && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update -y && \
    apt-get install -y golang-1.11 ssh git make

ARG SSH_KEY
RUN mkdir -p  ~/.ssh && \
    echo "$SSH_KEY" | tr -d '\r' > ~/.ssh/id_rsa && \
    chmod 700 ~/.ssh/id_rsa && \
    ssh-keyscan -H gitlab.com >> ~/.ssh/known_hosts && \
    git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/" && \
    mkdir /app
WORKDIR /app

COPY . /app

ENV GOPATH /usr/src/go
ENV GOROOT /usr/lib/go-1.11
ENV PATH $PATH:$GOROOT/bin

CMD tail -f /dev/null
