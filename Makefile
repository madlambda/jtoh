version?=$(shell git rev-list -1 HEAD)
cov=coverage.out
covhtml=coverage.html

all: build

.PHONY: build
build:
	go build -i  -ldflags "-X main.Version=${version}" .

.PHONY: lint
lint:
	@echo TODO

.PHONY: test
test:
	go test -race -coverprofile=$(cov) ./...

.PHONY: coverage
coverage: test
	go tool cover -html=$(cov) -o=$(covhtml)
