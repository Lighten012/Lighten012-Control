GO ?= go
GO_TAGS ?= remote exclude_graphdriver_btrfs containers_image_openpgp

.PHONY: build run vet test fmt clean

build:
	$(GO) build -tags "$(GO_TAGS)" -o bin/lighten012-control ./cmd/server

run:
	$(GO) run -tags "$(GO_TAGS)" ./cmd/server

vet:
	$(GO) vet -tags "$(GO_TAGS)" ./...

test:
	$(GO) test -tags "$(GO_TAGS)" ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf bin
