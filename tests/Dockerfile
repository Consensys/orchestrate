############################
# STEP 1 build executable binary
############################
FROM golang:1.14 AS builder

RUN useradd appuser && \
    mkdir /app
WORKDIR /app

# Use go mod with go 1.14
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/main -a -tags netgo -ldflags '-w -s -extldflags "-static"' ./tests/cmd

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/main /go/bin/main
COPY --from=builder /app/tests/features /features

# Use an unprivileged user.
USER appuser
EXPOSE 8080
ENTRYPOINT ["/go/bin/main"]
