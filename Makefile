.PHONY: build
build:
	go build -v ./cmd/portoEvents

.PHONY: run
run:
	go run ./cmd/portoEvents

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

#.DEFAULT_GOAL := run
.DEFAULT_GOAL := build

PROJECT_NAME="PortoEvents"
STDERR="log-stderr.txt"
#STDERR=./log/.$(PROJECT_NAME)-stderr.txt

#.DEFAULT_GOAL := generate