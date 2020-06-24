version?=$(shell git rev-list -1 HEAD)
cov=coverage.out
covhtml=coverage.html

all: build

.PHONY: build
build:
	go build -i  -ldflags "-X main.Version=${version}" .

.PHONY: lint
lint:
	@golangci-lint run ./... \
	    --enable=unparam --enable=unconvert --enable=dupl --enable=gofmt \
	    --enable=stylecheck --enable=scopelint --enable=nakedret --enable=misspell \
	    --enable=goconst --enable=dogsled --enable=bodyclose --enable=whitespace

.PHONY: test
test:
	go test -timeout 10s -race -coverprofile=$(cov) ./...

.PHONY: coverage
coverage: test
	go tool cover -html=$(cov) -o=$(covhtml)
