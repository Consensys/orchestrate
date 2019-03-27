# API-Context-Store

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **GRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

API-Context-Store is responsible to store Transaction execution context while transaction is being mined 

## Quick-Start

### Prerequisite

- Having ```docker``` and ```docker-compose``` installed
- Having Go 1.12 installed or upper

### Start the application

To quickly start the application

1. Start e2e environment

```sh
$ docker-compose -f e2e/docker-compose.yml
```

2. Migrate database

```sh
$ go run . migrate init
$ go run . migrate
```

3. Start worker

```sh
$ go run . run
```

### Configure Run

Application can be configured through flags or environment variables, you can run the ```help run``` command line

```sh
$ go run . help run
```

```
Run application

Usage:
  app run [flags]

Flags:
  -h, --help                   help for run
      --http-hostname string   Hostname to expose healthchecks and metrics.
                               Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --jaeger-host string     Jaeger host.
                               Environment variable: "JAEGER_HOST" (default "jaeger")
      --jaeger-port int        Jaeger port
                               Environment variable: "JAEGER_PORT" (default 5775)
      --jaeger-sampler float   Jaeger sampler
                               Environment variable: "JAEGER_SAMPLER" (default 0.5)

Global Flags:
      --db-database string   Target Database name
                             Environment variable: "DB_DATABASE" (default "postgres")
      --db-host string       Database host
                             Environment variable: "DB_HOST" (default "127.0.0.1")
      --db-password string   Database User password
                             Environment variable: "DB_PASSWORD" (default "postgres")
      --db-poolsize int      Maximum number of connections on database
                             Environment variable: "DB_POOLSIZE"
      --db-port int          Database port
                             Environment variable: "DB_PORT" (default 5432)
      --db-user string       Database User.
                             Environment variable: "DB_USER" (default "postgres")
      --log-format string    Log formatter (one of ["text" "json"]).
                             Environment variable: "LOG_FORMAT" (default "text")
      --log-level string     Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                             Environment variable: "LOG_LEVEL" (default "debug")
```

### Configure migrate
Migration can be configured through flags or environment variables, you can run the ```help migrate``` command line

```sh
$ go run . help migrate
```

```
Migrate database

Usage:
  app migrate [flags]
  app migrate [command]

Available Commands:
  down        Reverts last migration
  init        Initialize database
  reset       Reverts all migrations
  set-version Set database version
  up          Upgrade database
  version     Print current database version

Flags:
  -h, --help   help for migrate

Global Flags:
      --db-database string   Target Database name
                             Environment variable: "DB_DATABASE" (default "postgres")
      --db-host string       Database host
                             Environment variable: "DB_HOST" (default "127.0.0.1")
      --db-password string   Database User password
                             Environment variable: "DB_PASSWORD" (default "postgres")
      --db-poolsize int      Maximum number of connections on database
                             Environment variable: "DB_POOLSIZE"
      --db-port int          Database port
                             Environment variable: "DB_PORT" (default 5432)
      --db-user string       Database User.
                             Environment variable: "DB_USER" (default "postgres")
      --log-format string    Log formatter (one of ["text" "json"]).
                             Environment variable: "LOG_FORMAT" (default "text")
      --log-level string     Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                             Environment variable: "LOG_LEVEL" (default "debug")

Use "app migrate [command] --help" for more information about a command.
```