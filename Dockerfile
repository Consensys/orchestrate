ARG VERSION=nonroot

############################
# STEP 1 build executable Orchestrate binary
############################
FROM golang:1.16.9 AS builder

RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	ca-certificates upx-ucl

RUN useradd appuser && mkdir /app
WORKDIR /app

# Use go mod with go 1.15
ENV GO111MODULE=on
COPY go.mod go.sum ./
COPY LICENSE ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /bin/main -a -tags netgo -ldflags '-w -s -extldflags "-static"' .
RUN upx /bin/main

############################
# STEP 2 build a small image
############################
FROM gcr.io/distroless/static:$VERSION

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/main /go/bin/main
COPY --from=builder /app/LICENSE /

# Use an unprivileged user.
USER appuser
EXPOSE 8080
ENTRYPOINT ["/go/bin/main"]
