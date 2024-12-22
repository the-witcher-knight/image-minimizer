# Variables
PROJECT_NAME=image-minimize-go

# Shorten
DOCKER_BUILD_BIN=docker
COMPOSE_BIN=PROJECT_NAME=$(PROJECT_NAME) docker compose
COMPOSE_RUN := $(COMPOSE_BIN) run --rm --service-ports server

# Build
.PHONY: build-image-dev
build-image-dev:
	@${DOCKER_BUILD_BIN} build -f Dockerfile -t ${PROJECT_NAME}-local:latest .
	-${DOCKER_BUILD_BIN} images -q -f "dangling=true" | xargs docker rmi -f

run:
	@${COMPOSE_RUN} sh -c "go run ./cmd/serverd"