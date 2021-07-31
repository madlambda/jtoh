version?=$(shell git rev-list -1 HEAD)
cov=coverage.out
covhtml=coverage.html
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
	go test -timeout 10s -race -coverprofile=$(cov) ./...

bench: name?=.
bench:
	go test -bench=$(name) -benchmem -memprofile=mem.prof -cpuprofile cpu.prof .

bench/mem/analyze:
	go tool pprof mem.prof

bench/cpu/analyze:
	go tool pprof cpu.prof

.PHONY: coverage
coverage: test
	go tool cover -html=$(cov) -o=$(covhtml)

.PHONY: install
install:
	go install $(buildflags) ./cmd/jtoh

.PHONY: cleanup
cleanup:
	rm -f *.prof
	rm -f jtoh.test
	rm -f cmd/jtoh/jtoh
