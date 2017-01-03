.PHONY: compile build build_all fmt lint test itest vet bootstrap
	
SOURCE_FOLDER := .

BINARY_PATH ?= ./bin/configo

GOARCH ?= amd64

ifdef GOOS
BINARY_PATH :=$(BINARY_PATH).$(GOOS)-$(GOARCH)
endif

SPECS ?= spec/integration/**/**

default: build
	
build_all: vet fmt
	for GOOS in darwin linux windows; do \
		$(MAKE) compile GOOS=$$GOOS GOARCH=amd64 ; \
	done

compile:
	CGO_ENABLED=0 go build -i -v -ldflags '-s' -o $(BINARY_PATH) $(SOURCE_FOLDER)/

build: vet fmt compile
	
fmt:
	go fmt $(glide novendor)

vet:
	go vet $(glide novendor)

lint:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 golint

test:
	go test $(glide novendor)
	
itest:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	bats $(SPECS)

bootstrap:
	go get github.com/Masterminds/glide
	glide install
