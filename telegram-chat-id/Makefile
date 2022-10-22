PROJECT_NAME=telegram-chat-id
IMAGE_NAME=telegram-chat-id
export DOMAIN=http://localhost:8081
export FAASD_DOMAIN=openfaas.k3s.dansult.space
export OPENFAAS_PREFIX=docker.io/danielmichaels
export OPENFAAS_URL=https://${FAASD_DOMAIN}
export DOCKER_COMPOSE=true

.PHONY: help

default: help

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

## help: Print commands help.
.PHONY: help
help : Makefile
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## env: print environment variables (makefile sanity check)
.PHONY: env
env:
	env

## mksecrets: make secret files for development
.PHONY: mksecrets
mksecrets:
	mkdir -p /tmp/secrets/

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## dev/build: Build the image locally
.PHONY: dev/build
dev/build:
	make mksecrets
	@faas-cli build -f $(PROJECT_NAME).yml

## dev/run: Run the image locally
.PHONY: dev/run
dev/run:
#	make dev/build
	make mksecrets
	@faas-cli build -f $(PROJECT_NAME).yml
	@docker run --rm -it -p 8081:8080 danielmichaels/$(IMAGE_NAME)
## dev/reload: Hot-reloading locally using Air and Docker-compose
## credits: https://simonemms.com/blog/2020/08/12/live-reload-for-openfaas/
.PHONY: dev/reload
dev/reload:
	make mksecrets
	@docker-compose up --build
# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## prod/build: Build the image locally
.PHONY: prod/build
prod/build:
	faas-cli build -f $(PROJECT_NAME).yml

## prod/push: Push the image to dockerhub
.PHONY: prod/push
prod/push:
	@faas-cli push -f $(PROJECT_NAME).yml

## prod/deploy: Deploy to faasd
.PHONY: prod/deploy
prod/deploy:
	@faas-cli deploy -f $(PROJECT_NAME).yml -e DOMAIN=$(OPENFAAS_URL)

## prod/up: Run faas-cli up
.PHONY: prod/up
prod/up:
	@faas-cli up -f $(PROJECT_NAME).yml -e DOMAIN=$(OPENFAAS_URL)

## prod/logs: Get a functions logs
.PHONY: prod/logs
prod/logs:
	@faas-cli logs ${PROJECT_NAME}
