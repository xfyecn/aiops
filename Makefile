APP=aiops
VERSION=4.2
SHELL=/bin/bash

all: deps build

clean:
	@echo "--> cleaning..."
	@rm -rf build
	@go clean ./...

prereq:
	@mkdir -p build/{bin,tar}
	@go get -u github.com/Masterminds/glide

deps: prereq
	@glide install --strip-vendor

build: prereq
	@echo '--> building...'
	@go fmt ./...
	go build -o build/bin/${APP} ./cmd

package:
	@echo '--> packaging...'
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o build/bin/${APP} ./cmd
	@docker build -t rivernet/aiops:${VERSION} .

deploy: package
	@echo '--> deploying...'
	@docker push rivernet/aiops:${VERSION}
