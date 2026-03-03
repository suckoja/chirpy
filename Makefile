# Makefile for github.com/suckoja/chirpy
#
# Usage:
#   make build
#   make run
#   make test
#   make fmt
#   make lint
#   make clean

APP_NAME := chirpy
CMD_DIR  := ./cmd/$(APP_NAME)
BIN_DIR  := ./bin
BIN      := $(BIN_DIR)/$(APP_NAME)

GO      := go
GOFLAGS :=

.PHONY: help build run test tidy fmt vet lint clean

help:
	@echo "Targets:"
	@echo "  build  - build the binary to $(BIN)"
	@echo "  run    - run the app"
	@echo "  test   - run tests"
	@echo "  tidy   - go mod tidy"
	@echo "  fmt    - gofmt all Go files"
	@echo "  vet    - go vet"
	@echo "  lint   - fmt + vet + test"
	@echo "  clean  - remove build artifacts"

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN) $(CMD_DIR)

run:
	$(GO) run $(GOFLAGS) $(CMD_DIR)

test:
	$(GO) test $(GOFLAGS) ./...

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint: fmt vet test

clean:
	@rm -rf $(BIN_DIR)