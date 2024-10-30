# Makefile for the Taco project

# Variables
SRC_DIR := ./src
OUTPUT := taco
INSTALL_DIR := /usr/local/bin

# Default run target (runs the project without any arguments)
.PHONY: run
run:
	go run $(SRC_DIR)

# Run with custom output file
.PHONY: run-output
run-output:
	go run $(SRC_DIR) -output=$(OUTPUT).txt

# Build the project (creates the taco executable)
.PHONY: build
build:
	go build -o $(OUTPUT) $(SRC_DIR)

# Clean the project (removes the taco executable)
.PHONY: clean
clean:
	rm -f $(OUTPUT)

# Run tests (if you add tests in the future)
.PHONY: test
test:
	go test -v ./...

# Install the taco executable to /usr/local/bin
.PHONY: install
install: build
	sudo install -m 755 $(OUTPUT) $(INSTALL_DIR)

# Uninstall the taco executable from /usr/local/bin
.PHONY: uninstall
uninstall:
	rm -f $(INSTALL_DIR)/$(OUTPUT)
