PROJECT:=http-server

BINARY = http-server
SOURCE_DIR = ./

appname := http-server

sources := $(wildcard *.go)

build = CGO_ENABLE=0 GOOS=$(1) GOARCH=$(2) go build -o build/$(appname)$(3) $(SOURCE_DIR)

.PHONY: build

all: clean build

clean:
	rm -rf build

build:
	go build -o build/$(BINARY) $(SOURCE_DIR)

linux: $(sources)
	$(call build,linux,amd64,)

windows: $(sources)
	$(call build,windows,amd64,.exe)

darwin: $(sources)
	$(call build,darwin,amd64,)
