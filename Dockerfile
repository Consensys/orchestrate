############################
# STEP 1 build executable binary
############################
FROM golang:1.11 as builder

ARG SSH_KEY
RUN useradd appuser && \
    mkdir -p  ~/.ssh && \
    echo "$SSH_KEY" | tr -d '\r' > ~/.ssh/id_rsa && \
    chmod 700 ~/.ssh/id_rsa && \
    ssh-keyscan -H gitlab.com >> ~/.ssh/known_hosts && \
    git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/" && \
    mkdir /app
WORKDIR /app

# Use go mod with go 1.11
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/main -a -tags netgo -ldflags '-w -s -extldflags "-static"' .


############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/main /go/bin/main

# Label image
LABEL org.label-schema.schema-version="1.0.0-rc1"

# Use an unprivileged user.
USER appuser
EXPOSE 8080
ENTRYPOINT ["/go/bin/main"]
CMD ["run"]
