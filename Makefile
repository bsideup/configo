.PHONY: compile build build_all fmt lint test itest vet godep_save godep_restore
	
SOURCE_FOLDER := .

BINARY_PATH ?= ./bin/configo

ifdef GOOS
BINARY_PATH :=$(BINARY_PATH).$(GOOS)-$(GOARCH)
endif

export GO15VENDOREXPERIMENT=1

default: build
	
build_all: vet fmt
	for GOOS in darwin linux windows; do \
		for GOARCH in 386 amd64; do \
			$(MAKE) compile GOOS=$$GOOS GOARCH=$$GOARCH ; \
		done \
	done

compile:
	CGO_ENABLED=0 go build -ldflags '-s' -o $(BINARY_PATH) $(SOURCE_FOLDER)/

build: vet fmt compile
	
fmt:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 go fmt

vet:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 go vet

lint:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 golint

test:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 go test
	
itest:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	bats spec/integration/**

godep_save:
	go get github.com/tools/godep
	godep save $(SOURCE_FOLDER)/...

godep_restore:
	go get github.com/tools/godep
	godep restore -v
