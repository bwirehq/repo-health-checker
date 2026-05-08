.PHONY: build test vet lint fmt check

build:
	go build -o bin/repo-health ./cmd/repo-health

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

check: test vet lint
