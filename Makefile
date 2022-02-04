BIN := "./cmd/bin"

build:
	go build -v -o $(BIN) ./cmd

run: build
	$(BIN)

test:
	go test -v ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.39.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run test lint