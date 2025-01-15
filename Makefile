.PHONY: build install server

PROJECT_NAME=$(shell basename $(CURDIR))

## build: build the application
build:
	export GO111MODULE="auto"; \
	go mod download; \
	go mod vendor; \
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o main main.go

## install: fetches go modules
install:
	export GO111MODULE="on"; \
	go mod tidy; \
	go mod download \

## server: runs the server with -race
server:
	export GO111MODULE="on"; \
	go run main.go

## modd: Monitors the directory, and recompiles your app every time a file changes. Also runs tests.
## (To install modd, run: go get github.com/cortesi/modd/cmd/modd)
modd:
	modd
