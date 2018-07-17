BINARY_NAME=aws-env
TARGET ?= darwin
ARCH ?= amd64
EXT ?= ""
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: test

generate:
	@echo "== Go Generate =="
	go generate ./...

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)
	go vet -v ./...
	go test -race -v ./...

run: test
	@echo "== Run =="
	go run cmd/main.go

build: test
	@echo "== Build =="
	go build -o $(BINARY_NAME) -v cmd/main.go

release: test
	@echo "== Release build =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o $(BINARY_NAME)-$(TARGET)-$(ARCH)$(EXT) -v cmd/main.go

clean:
	@echo "== Cleaning =="
	rm $(BINARY_NAME)* || true

.PHONY: default generate test run build release clean
