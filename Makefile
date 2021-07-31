version?=$(shell git rev-list -1 HEAD)
cov=coverage.out
covhtml=coverage.html
buildflags=-ldflags "-X main.Version=${version}"

all: build test lint

.PHONY: build
build:
	go build -i  $(buildflags) -o ./cmd/jtoh/jtoh ./cmd/jtoh

.PHONY: lint
lint:
	golangci-lint run ./... \
	    --enable=unparam --enable=unconvert --enable=dupl --enable=gofmt \
	    --enable=stylecheck --enable=scopelint --enable=nakedret --enable=misspell \
	    --enable=goconst --enable=dogsled --enable=bodyclose --enable=whitespace --enable=golint

.PHONY: test
test:
	go test -timeout 10s -race -coverprofile=$(cov) ./...

bench: name?=.
bench:
	go test -bench=$(name) -benchmem -memprofile=memory.p .

bench/mem/analyze:
	go tool pprof memory.p

.PHONY: coverage
coverage: test
	go tool cover -html=$(cov) -o=$(covhtml)

.PHONY: install
install:
	go install $(buildflags) ./cmd/jtoh
