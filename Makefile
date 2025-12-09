.PHONY: build ui test

build:
	go build

ui:
	go run . -v ui

test:
	go test -tags=test ./...
