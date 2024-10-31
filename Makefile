# Makefile for the Taco project

# Variables
SRC_DIR := ./src
OUTPUT := taco
INSTALL_DIR := /usr/local/bin
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

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

# Clean the project (removes the taco executable and coverage files)
.PHONY: clean
clean:
	rm -f $(OUTPUT) $(COVERAGE_OUT) $(COVERAGE_HTML)

# Set the local Git config to use the custom hooks directory (./.git-hooks)
.PHONY: setup-git-hooks
setup-git-hooks:
	@echo "Setting up git hooks..."
	git config core.hooksPath .git-hooks
	chmod +x .git-hooks/*
	@echo "Git hooks have been configured successfully."

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage report
.PHONY: test-coverage
test-coverage:
	go test -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -func=$(COVERAGE_OUT)

# Generate HTML coverage report
.PHONY: coverage-html
coverage-html: test-coverage
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

# Install the taco executable to /usr/local/bin
.PHONY: install
install: build
	sudo install -m 755 $(OUTPUT) $(INSTALL_DIR)

# Uninstall the taco executable from /usr/local/bin
.PHONY: uninstall
uninstall:
	rm -f $(INSTALL_DIR)/$(OUTPUT)
