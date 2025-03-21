.PHONY: build
build:
	go build -v -o rpgapp ./cmd/app

.PHONY: test
test:
	go test -v -race -timeout 30s ./...
.DEFAULT_GOAL := build