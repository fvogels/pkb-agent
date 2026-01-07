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
	echo "Don't forget to enable profiling in the code and run it before profiling!"
	go tool pprof pkb-agent.exe ./profile.txt