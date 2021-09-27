default: build
NAME:=gitlab-multirepo-deployer

ifeq ($(VERSION),)
  VERSION_TAG=$(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null || git describe --abbrev=0 --all)-$(shell date +"%Y%m%d%H%M%S")
else
  VERSION_TAG=$(VERSION)-$(shell date +"%Y%m%d%H%M%S")
endif

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


.PHONY: help
help:
	@cat docs/developer-guide/make-targets.md

.PHONY: release
release: linux darwin compress

.PHONY: build
build:
	go build -o ./.bin/$(NAME) -ldflags "-X \"main.version=$(VERSION_TAG)\""  main.go


.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o .bin/$(NAME)_linux-amd64
	cp .bin/$(NAME)_linux-amd64 .bin/$(NAME)
	GOOS=linux GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o .bin/$(NAME)_linux-arm64

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\""  -o .bin/$(NAME)_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\"" -o .bin/$(NAME)_darwin-arm64

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X \"main.version=$(VERSION_TAG)\""  -o .bin/$(NAME).exe

.PHONY: compress
compress:
	upx -5 ./.bin/*

.PHONY: install
install:
	cp ./.bin/$(NAME) /usr/local/bin/

.PHONY: lint
lint: build
	golangci-lint run --verbose --print-resources-usage