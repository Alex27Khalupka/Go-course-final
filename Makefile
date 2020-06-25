.PHONY: build
build:
	go build -v ./cmd/apiserver
migrations:
	createdb tasks
	migrate -path migrations -database "postgres://localhost:5432/tasks?sslmode=disable&user=postgres&password=postgres" up
.DEFAULT_GOAL := build