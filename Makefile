.PHONY: build
build:
	go build -v ./cmd/apiserver
.PHONY: migrations
migrations:
	migrate -path migrations -database $$POSTGRES_URI up
.DEFAULT_GOAL := build