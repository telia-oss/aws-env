BINARY  = aws-env
TARGET ?= darwin
ARCH   ?= amd64
EXT    ?= ""

TRAVIS_TAG ?= $(shell git describe --tags --candidates=1 --dirty 2>/dev/null || echo "dev")
SRC         = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

export GO111MODULE=on

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

e2e: test
	@echo "== Integration =="
	go test -race -v ./... -tags=e2e

build: test
	@echo "== Build =="
	go build -o $(BINARY) -v cmd/main.go

release: test
	@echo "== Release build =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -ldflags="-s -w -X=main.version=$(TRAVIS_TAG)" -o $(BINARY)-$(TARGET)-$(ARCH)$(EXT) -v cmd/main.go

clean:
	@echo "== Cleaning =="
	rm $(BINARY)* || true

.PHONY: default generate test run build release clean
