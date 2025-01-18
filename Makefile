.PHONY: build install server

PROJECT_NAME=$(shell basename $(CURDIR))

APPLICATION_NAME := traefik-etcd-advertiser
VERSION := $(shell git describe --tags HEAD 2>/dev/null || echo $(DRONE_TAG))
# PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 windows/amd64
PLATFORMS := darwin/arm64 linux/amd64

## build: build the application
build: install
	export GO111MODULE="auto"; \
		go mod download; \
		go mod vendor; \
		CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o main main.go


release: install-vendor
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		os=$$GOOS; \
		if [ "$$os" = "darwin" ]; then \
			os="macOS"; \
		fi; \
		output_name="$(APPLICATION_NAME)-$(VERSION)-$$os-$$GOARCH"; \
		if [ "$$os" = "windows" ]; then \
			output_name+=".exe"; \
		fi; \
		echo "Building release/$$output_name..."; \
		env GO11MODULE="auto" CGO_ENABLED=0 \
			GOOS=$$GOOS GOARCH=$$GOARCH go build \
			-ldflags "-w -s -X main.Version=$(VERSION)" \
			-o release/$$output_name; \
		if [ $$? -ne 0 ]; then \
			echo 'An error has occurred! Aborting.'; \
			exit 1; \
		fi; \
		cd release > /dev/null; \
		if [ "$$os" = "windows" ]; then \
			zip "$$output_name.zip" "$$output_name"; \
			rm "$$output_name"; \
		else \
			chmod a+x "$$output_name"; \
			tar -czf "$$output_name.tar" "$$output_name"; \
			rm "$$output_name"; \
		fi; \
		cd ..; \
	done

## install: fetches go modules
install:
	export go111module="on"; \
	go mod tidy; \
	go mod download \

install-vendor: install
	export go111module="on"; \
	go mod vendor

## server: runs the server with -race
server:
	export GO111MODULE="on"; \
	go run main.go

## modd: Monitors the directory, and recompiles your app every time a file changes. Also runs tests.
## (To install modd, run: go get github.com/cortesi/modd/cmd/modd)
modd:
	modd
