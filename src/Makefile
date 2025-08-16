# Generic Go Makefile

APP_NAME := go_grep
PKG := ./...
BIN := $(APP_NAME)

.PHONY: all build run fmt vet

all: build

## Build binary
build:
	go build -o $(BIN) .

## Run the binary (with optional args: make run ARGS="foo bar")
run: build
	./$(BIN) $(ARGS)

## Format code
fmt:
	go fmt $(PKG)

## Vet code (find suspicious constructs)
vet:
	go vet $(PKG)
