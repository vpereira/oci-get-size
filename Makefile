PROJECT_NAME := oci-get-size

# Build target
build:
	@echo "Building the project..."
	go build -o $(PROJECT_NAME) main.go

# Test target
test:
	@echo "Running tests..."
	go test -v ./...

# Format target
format:
	@echo "Formatting the project..."
	go fmt ./...

# Default target
.PHONY: all build test format

all: build

