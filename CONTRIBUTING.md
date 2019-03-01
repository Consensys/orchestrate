# Contributing Guidelines

:dolphin: Thanks for being interested in contributing to Core-Stack! People will :heart: you for that. :thumbsup:

## Create dev environment

### Prerequisite

Core-Stack requires 

- Go 1.11
- ```docker``` & ```docker-compose```
- ```minikube``` (optional)

### Important note

If you are developing for the first time on the project you might run into trouble when installing dependencies. This is due to the fact that we are using ```go``` dependency from private GitLab repositories.

To avoid issues we recommend you setting an SSH Key on your GitLab account and then run locally

```sh
$ git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
```

### Installation

To contribute

1. Clone project locally

```sh
$ git clone git@gitlab.com:ConsenSys/client/fr/core-stack/<project_path>.git
$ cd <project_name>
```

2. Test everything is okay

```sh
$ make run-coverage
```

### Starting e2e env

1. Start e2e environment by running

```sh
$ docker-compose -f e2e/docker-compose.yml up
```

2. Start worker locally by running

```sh
$ go run . run
```

3. Produce messages (from another terminal)

```sh
$ go run e2e/producer/main.go
```

## Project Structure

```text
.
├── app/
│   ├── infra/                      # infrastructure resources (connectors to Ethereum clients, Kafka, Databases...)
│   ├── worker/                     # worker responsible to treat consumed messages on Kafka topic
|   │   ├── handlers/               # handlers to be register on worker 
│   │   └── worker.go               # create worker
│   ├── app.go                      # main application object
│   ├── consumer_group.go           # kafka consumer group
│   └── server.go                   # server to expose metrics, liveness, readyness check...
├── cmd/                            # command line interface (we use Cobra)
│   ├── root.go                     # root command from which all CLI commands inherit
│   └── run.go                      # run command to start the application
├── e2e/                            # end 2 end utilities to test application
│   ├── consumer/                   # define a consumer listening on worker output kafka topic
│   ├── producer/                   # define a producer producing on worker input kafka topic
│   └── docker-compose.yml          # docker-compose file to start a local Kafka
├── infra/                          # infrastructure elements (this will be soon moved in an other repo)
├── scripts/                        # facility scripts for dev & CI/CD
├── .dockerignore                   # list files to be ignored when building docker image
├── .gitignore                      # list untracked files by git
├── .gitlab-ci.yml                  # CI/CD pipeline script
├── CHANGELOG.md                    # indicate list of changes from a release to another
├── CONTRIBUTING.md                 # contributing guidelines                
├── docker-compose.yml              # docker-compose used to build docker image in CI/CD
├── Dockerfile                      # dockerfile for the application
├── go.mod                          # go dependencies list
├── go.sum                          # go dependencies locked list
├── LICENSE                         # LICENSE
├── main.go                         # main file to start the application
├── Makefile                        # utility make command
└── README.md                       # README
```

## Git Branching Strategy

![alt git-branching-strategy](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/blob/master/diagrams/Git_Branching_Strategy.png)

For more details about git branching strategy refer to [Git branching Strategy](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/blob/master/doc/git_branching_strategy.md)


## Build

To be able to build a project which uses private repositories

```bash
# Run this command
export SSH_KEY=`cat ~/.ssh/id_rsa`

# Or add is to your .bashrc
echo 'export SSH_KEY=`cat ~/.ssh/id_rsa`' >> ~/.bashrc
```

*Note: if your private repository ssh key is not id_rsa, replace it in the above command.*