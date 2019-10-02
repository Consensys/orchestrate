############################
# STEP 1 build executable binary
############################
FROM golang:1.13 AS builder

ARG GITLAB_USER
ARG GITLAB_TOKEN

RUN git config --global --add url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "git@gitlab.com:" && \
    git config --global --add url."https://${GITLAB_USER}:${GITLAB_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/" && \
    useradd appuser && \
    mkdir /app
WORKDIR /app

# Use go mod with go 1.13
ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.com/ConsenSys/client/fr/core-stack
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/main -a -tags netgo -ldflags '-linkmode external -w -s' .

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/main /go/bin/main

# Use an unprivileged user.
USER appuser
EXPOSE 8080
ENTRYPOINT ["/go/bin/main"]
CMD ["run"]
