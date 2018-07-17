BINARY_NAME=aws-env
TARGET ?= linux
ARCH ?= amd64
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
DIR=$(shell pwd)

default: test

generate:
	@echo "== Go Generate =="
	go generate ./...

run: test
	@echo "== Run =="
	go run cmd/main.go

build: test
	@echo "== Build =="
	go build -o $(BINARY_NAME) -v cmd/main.go

clean:
	@echo "== Cleaning =="
	rm $(BINARY_NAME) || true
	rm concourse-github-lambda.zip || true

release:
	@echo "== Release build =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o $(BINARY_NAME) -v cmd/main.go

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)
	go vet -v ./...
	go test -race -v ./...

.PHONY: default build test release test-code generate
