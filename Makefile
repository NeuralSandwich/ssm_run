all: build

# Makefile for building and deploying the ssm_run application

# To build and run the docker container locally, run:
# $ make

#### VARIABLES ####
USERNAME = davyj0nes
APP_NAME = ssm_run
PROJECT ?= github.com/davyj0nes/ssm_run

IMAGE_VERSION ?= latest

GO_VERSION ?= 1.10

RELEASE = 0.0.2
COMMIT = $(shell git rev-parse HEAD | cut -c 1-6)
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

BUILD_PREFIX = CGO_ENABLED=0 GOOS=linux
BUILD_FLAGS = -a -tags netgo --installsuffix netgo
LDFLAGS = -ldflags "-s -w -X ${PROJECT}/cmd.Release=${RELEASE} -X ${PROJECT}/cmd.Commit=${COMMIT} -X ${PROJECT}/cmd.BuildTime=${BUILD_TIME}"
DOCKER_GO_BUILD = docker run --rm -v "$(GOPATH)":/go -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION}
GO_BUILD_STATIC = $(BUILD_PREFIX) go build $(BUILD_FLAGS) $(LDFLAGS)
GO_BUILD_OSX = GOOS=darwin GOARCH=amd64 go build $(LDFLAGS)
GO_BUILD_WIN = GOOS=windows GOARCH=amd64 go build $(LDFLAGS)W

DOCKER_RUN_CMD = docker run -it --rm -v ${APP_NAME}:/app/.tasks --name ${APP_NAME} ${USERNAME}/${APP_NAME}:${IMAGE_VERSION} "\$$@"

.PHONY: compile build install test clean

#### COMMANDS ####
compile:
	@mkdir -p releases/${RELEASE}
	$(call blue, "# Compiling Static Golang App...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_STATIC} -o ${APP_NAME}_static'
	$(call blue, "# Compiling OSX Golang App...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_OSX} -o releases/${RELEASE}/${APP_NAME}_osx'
	$(call blue, "# Compiling Windows Golang App...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_WIN} -o releases/${RELEASE}/${APP_NAME}.exe'

build: compile
	$(call blue, "# Building Docker Image...")
	@docker build --no-cache --label APP_VERSION=${RELEASE} --label BUILT_ON=${BUILD_TIME} --label GIT_HASH=${COMMIT} -t ${USERNAME}/${APP_NAME}:${IMAGE_VERSION} .
	@docker volume create ${APP_NAME}
	@$(MAKE) clean

install_binary: compile
	$(call blue, "# Installing binary...)
	ifeq($(shell uname -a), Darwin)
	AWS_DEFAULT_REGION=eu-west-1 AWS_PROFILE=pt-m aws ec2 describe-instances --filters dns-name=ip-10-65-128-74.eu-west-1.compute.internal
		@cp releases/${RELEASE}/${APP_NAME}_osx $(GOPATH/bin/${APP_NAME})
	endif

install: build
	$(call blue, "# Installing Docker Image Locally...")
	@rm -f $(HOME)/bin/ssm_rund
	@echo "#!/bin/bash" >> $(HOME)/bin/ssm_rund
	@echo "set -e" >> $(HOME)/bin/ssm_rund
	@echo ${DOCKER_RUN_CMD} >> $(HOME)/bin/ssm_rund
	@chmod +x $(HOME)/bin/ssm_rund

test:
	$(call blue, "# Linting Code...")
	@golint -min_confidence=0.3 ./...
	$(call blue, "# Running Tests...")
	@docker run --rm -it -v "$(GOPATH):/go" -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION} sh -c 'go test -v' 

clean: 
	@rm -f ${APP_NAME}_static

#### FUNCTIONS ####
define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef
