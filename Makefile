.PHONY: build clean tool lint help

all: build

build:
	@go build -v .

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf go-gin-example
	go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
dev:
	@echo "Starting development server with hot reload..."
	@air

dev-fast:
	@echo "Starting development server with Fresh (fastest)..."
	@fresh

dev-reflex:
	@echo "Starting development server with Reflex..."
	@reflex -c reflex.conf

install-tools:
	@echo "Installing all hot reload tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/gravityblast/fresh@latest
	@go install github.com/cespare/reflex@latest

install-air:
	@echo "Installing Air..."
	@go install github.com/cosmtrek/air@latest