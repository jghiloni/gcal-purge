GOCMD := go

GOBUILD    := $(GOCMD) build
GOGET      := $(GOCMD) get
GOTEST     := $(GOCMD) test

GOOS   := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)

OUT ?= build

.PHONY: all-os all clean test

all: clean build test
all-os: clean build-all test
build: gcal-purge
build-all: gcal-purge-all

gcal-purge: $(OUT)/gcal-purge-$(GOOS)

gcal-purge-all: $(OUT)/gcal-purge-linux $(OUT)/gcal-purge-darwin

$(OUT)/gcal-purge-linux: clean
	GOOS=linux $(GOBUILD) -o $(OUT)/gcal-purge-linux .

$(OUT)/gcal-purge-darwin: clean
	GOOS=darwin $(GOBUILD) -o $(OUT)/gcal-purge-darwin .

test:
	$(GOTEST) -race -covermode=atomic -cover ./...

clean:
	rm -f $(OUT)/*
