FROM golang:1.12

ARG GITLAB_USER
ARG GITLAB_TOKEN

RUN apt-get update -y && \
    apt-get install -y jq

RUN git config --global --add url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "git@gitlab.com:" && \
    git config --global --add url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/" && \
    useradd appuser && \
    mkdir /app
WORKDIR /app

COPY . /app

CMD tail -f /dev/null