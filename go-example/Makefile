PROJECT_NAME=go-example
IMAGE_NAME=go-faas-routes-example
export DOMAIN=http://localhost:8081
export FAASD_DOMAIN=openfaas.k3s.dansult.space
export OPENFAAS_PREFIX=docker.io/danielmichaels
export OPENFAAS_URL=https://${FAASD_DOMAIN}

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

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## dev/build: Build the image locally
.PHONY: dev/build
dev/build:
	@faas-cli build -f $(PROJECT_NAME).yml

## dev/run: Run the image locally
.PHONY: dev/run
dev/run:
#	make dev/build
	@docker run --rm -it -p 8081:8080 danielmichaels/$(IMAGE_NAME)

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
