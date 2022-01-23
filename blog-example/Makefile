PROJECT_NAME=blog-example
export ORIGINS=http://localhost:1313,https://danielms.site
export FAASD_DOMAIN=faasd.ptco.rocks
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
	@faas-cli build -f stack.yml

## dev/run: Run the image locally
.PHONY: dev/run
dev/run:
	@docker run -v /tmp/secrets/:/var/openfaas/secrets --rm -it -p 8080:8080 danielmichaels/$(PROJECT_NAME)

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## prod/build: Build the image locally
.PHONY: prod/build
prod/build:
	faas-cli build -f stack.yml

## prod/push: Push the image to dockerhub
.PHONY: prod/push
prod/push:
	@faas-cli push -f stack.yml

## prod/deploy: Deploy to faasd
.PHONY: prod/deploy
prod/deploy:
	@faas-cli deploy -f stack.yml

## prod/up: Run faas-cli up
.PHONY: prod/up
prod/up:
	@faas-cli up -f stack.yml

## prod/logs: Get a functions logs
.PHONY: prod/logs
prod/logs:
	@faas-cli logs ${PROJECT_NAME}
