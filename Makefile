GO := $(shell which go)

ifeq ($(GO),)
$(error "Go is not installed or not in PATH")
endif

.PHONY: build run test clean

build:
	@mkdir -p bin
	@$(GO) build -o bin/fs
	@chmod +x bin/fs

run: build
	@./bin/fs

test:
	@$(GO) test -v ./...

clean:
	@rm -rf bin
