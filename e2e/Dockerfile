FROM ubuntu:16.04

RUN apt-get update -y && \
    apt-get install -y software-properties-common jq curl && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update -y && \
    apt-get install -y golang-1.11 ssh git make

ARG GITLAB_USER
ARG GITLAB_TOKEN

RUN git config --global --add url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "git@gitlab.com:" && \
    git config --global --add url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/" && \
    useradd appuser && \
    mkdir /app
WORKDIR /app

COPY . /app

ENV GOPATH /usr/src/go
ENV GOROOT /usr/lib/go-1.11
ENV PATH $PATH:$GOROOT/bin

CMD tail -f /dev/null
