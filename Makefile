APP_NAME := OmniShard
SOURCE := .
BUILD_DIR := bin
GIT_VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X 'main.version=$(GIT_VERSION)' -X 'main.buildTime=$(BUILD_TIME)'"
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all build start install clean docker-build test help

all: build

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	@echo "Building $(APP_NAME) application..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SOURCE)
	@echo "Build done: $(APP_NAME)"

start: build
	@echo "Build and run $(APP_NAME) application..."
	./$(BUILD_DIR)/$(APP_NAME)

install: build
	@echo "Build and install $(APP_NAME) application..."
	cp ./$(BUILD_DIR)/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	@echo "Installation done"

clean: 
	@echo "Clean old version Go application..."
	rm -rf $(BUILD_DIR)/$(APP_NAME)
	@echo "Clean done"

docker-build: # WIP

test:
	@echo "Run Go tests..."
	go test -v ./...
	@echo "Tests done"

help:
	@echo "Usage lists:"
	@echo " make                      - Build the application (default)."
	@echo " make build                - Build the application."
	@echo " make start                - Build and run the application."
	@echo " make install              - Build and install the application."
	@echo " make clean                - Clean old version application."
	@echo " make docker-build         - (Unavaliable) Build the Docker image."
	@echo " make test                 - Run Go tests."
	@echo " make help                 - Show help messages."
