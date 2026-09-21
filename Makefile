GO ?= go

.PHONY: build run vet test fmt clean

build:
	$(GO) build -o bin/lighten012-control ./cmd/server

run:
	$(GO) run ./cmd/server

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf bin
