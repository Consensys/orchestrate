# Envelope Store API

  **Note**: Envelope store was previously called *Context Store*.

## Goal

API-Context-Store is responsible for storing Transaction execution context while it's being mined.

## Quick-Start

### Prerequisites

- Having ```docker``` and ```docker-compose``` installed;
- Having Go 1.12 installed or upper.

### Start the application

To quickly start the application:

**1. Start e2e environment**

```sh
$ docker-compose -f e2e/docker-compose.yml up
```

**2. Migrate database**

```sh
$ go run . migrate init
$ go run . migrate
```

**3. Start worker**

```sh
$ go run . run
```
