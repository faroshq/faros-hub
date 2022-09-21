REPO ?= quay.io/faroshq/kcp-service-examle
TAG_NAME ?= $(shell git describe --tags --abbrev=0)

run-server:
	@echo "Starting server..."
	@go run ./cmd/server

build-server:
	@echo "Building server..."
	@go build -o ./bin/server ./cmd/server


.PHONY: image-server
image-server:
	docker build -t ${REPO}:${TAG_NAME} -f dockerfiles/server/Dockerfile \
	--build-arg version=${TAG_NAME} .
