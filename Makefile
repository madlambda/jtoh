# Set default shell to bash
SHELL := /bin/bash -o pipefail -o errexit -o nounset

COVERAGE_REPORT ?= coverage.txt

version?=$(shell git rev-list -1 HEAD)
buildflags=-ldflags "-X main.Version=${version}"
golangci_lint_version=1.41.1

all: build test lint

.PHONY: build
build:
	go build -i  $(buildflags) -o ./cmd/jtoh/jtoh ./cmd/jtoh

.PHONY: lint
lint:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v$(golangci_lint_version)  golangci-lint run ./...

.PHONY: test
test:
	go test -timeout 10s -race ./...

.PHONY: test/fuzz
test/fuzz:
	go test -fuzz=FuzzJTOH

.PHONY: bench
bench: name?=.
bench:
	go test -bench=$(name) -benchmem -memprofile=mem.prof -cpuprofile cpu.prof .

.PHONY: time/bench
time/bench:
	time -v go test -bench=. .

.PHONY: bench/mem/analyze
bench/mem/analyze:
	go tool pprof mem.prof

.PHONY: bench/cpu/analyze
bench/cpu/analyze:
	go tool pprof cpu.prof

.PHONY: coverage
coverage: 
	go test -coverprofile=$(COVERAGE_REPORT) ./...

.PHONY: coverage/show
coverage/show: coverage
	go tool cover -html=$(COVERAGE_REPORT)

.PHONY: install
install:
	go install $(buildflags) ./cmd/jtoh

.PHONY: cleanup
cleanup:
	rm -f *.prof
	rm -f jtoh.test
	rm -f cmd/jtoh/jtoh
