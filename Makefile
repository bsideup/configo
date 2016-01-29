.PHONY: compile build build_all fmt lint test itest vet godep_save godep_restore
	
SOURCE_FOLDER := .

BINARY_PATH ?= ./bin/configo

ifdef GOOS
BINARY_PATH :=$(BINARY_PATH).$(GOOS)-$(GOARCH)
endif

# We have to specify them manually because of GO15VENDOREXPERIMENT bug (vendor folder not excluded)
PACKAGES := $(SOURCE_FOLDER) $(SOURCE_FOLDER)/sources/... $(SOURCE_FOLDER)/parsers/... $(SOURCE_FOLDER)/flatmap/...

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
	go fmt $(PACKAGES)

vet:
	go vet $(PACKAGES)

lint:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 golint

test:
	go test $(PACKAGES)
	
itest:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	bats spec/integration/**

godep_save:
	go get github.com/tools/godep
	godep save $(PACKAGES)

godep_restore:
	go get github.com/tools/godep
	godep restore -v
