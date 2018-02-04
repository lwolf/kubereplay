NAME := kube-replay
AUTHOR=lwolf
VERSION ?= 0.0.1
REGISTRY ?= quay.io
GIT_SHA=$(shell git --no-pager describe --always --dirty)
BUILD_TIME=$(shell date '+%s')
LFLAGS ?= -X main.gitsha=${GIT_SHA} -X main.compiled=${BUILD_TIME}
ROOT_DIR=${PWD}
GOVERSION ?= 1.9.3
HARDWARE=$(shell uname -m)

.PHONY: changelog build build-controller build-apiserver docker static release install_deps

golang:
	@echo "--> Go Version"
	@go version

build-info:
	@echo "building ${GIT_SHA} from ${BUILD_TIME}"

install_deps:
	dep ensure

build: golang build-info
	@echo "--> Compiling the project"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/${NAME}-controller cmd/controller/main.go
	go build -ldflags "${LFLAGS}" -o bin/${NAME}-apiserver cmd/apiserver/main.go

build-image: bin/linux/$(OPERATOR_NAME)
	docker build . -t $(IMAGE):$(VERSION)

build-controller: golang build-info
	@echo "--> Compiling the project"
	@mkdir -p bin
	go build -o bin/${NAME}-controller cmd/controller/main.go

build-apiserver: golang build-info
	@echo "--> Compiling the project"
	@mkdir -p bin
	go build -o bin/${NAME}-apiserver cmd/apiserver/main.go

clean:
	rm -rf ./bin 2>/dev/null
	rm -rf ./release 2>/dev/null

format:
	@echo "--> Running go fmt"
	@gofmt -s -w ./