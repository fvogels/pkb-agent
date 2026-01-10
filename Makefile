.PHONY: build ui test

build:
	go build

ui:
	go run . -v ui

test:
	go test -tags=test ./...

bench:
	go test -tags=test -bench=. ./...

profile:
	echo "Don't forget to run the application with profiling enabled (-b)!"
	go tool pprof pkb-agent.exe ./pkb-agent.prof