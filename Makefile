.PHONY: build
build:
	go build -v -o rpgapp ./cmd/app

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: migrate
migrate:
	migrate -path migrations/stories -database "postgres://postgres:world555@localhost:5431/rpg?sslmode=disable" down
	migrate -path migrations/stories -database "postgres://postgres:world555@localhost:5431/rpg?sslmode=disable" up
.DEFAULT_GOAL := build