# Contributing Guidelines

:dolphin: Thanks for being interested in contributing to Orchestrate! People will :heart: you for that. :thumbsup:

## Create dev environment

### Prerequisite

Orchestrate requires 

- Go 1.16
- ```docker``` & ```docker-compose```

### Installation

To contribute

1. Clone project locally

```sh
$ git clone git@github.com:ConsenSys/orchestrate.git <project_name> 
$ cd <project_name>
```

### Running local development environment

1. Start Orchestrate development

```sh
$ make dev
```

To start different test network you can run:

1. For geth client
```
$ make geth
```

2. For Besu client
```
$ make besu
```

3. For go-quorum client
```
$ make go-quorum
```

### Testing

1. Run linting checks
```sh
$ make lint
```

2. Run unit tests
```sh
$ make coverage
```

3. Run integration tests
```sh
$ make run-integration
```

4. Run end to end test suite

```sh
$ cp .env.ci .env
$ make e2e
```

## Git Branching Strategy

![alt git-branching-strategy](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/raw/master/diagrams/Git_Branching_Strategy.png)

For more details about git branching strategy refer to [Git branching Strategy](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/raw/master/diagrams/Git_Branching_Strategy.png)
