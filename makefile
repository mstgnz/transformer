.PHONY: live
.DEFAULT_GOAL:= run

run:
	clear && go build -o /tmp/build ./example && /tmp/build

live:
	find . -type f \( -name '*.go' \) | entr -r sh -c 'go build -o /tmp/build ./example && clear && /tmp/build'

# Test commands
.PHONY: test test-json test-xml test-yaml test-node test-benchmark test-cover test-bench

# Run all tests
test:
	@echo "Running all tests..."
	@go test ./...

# Run JSON tests
test-json:
	@echo "Running JSON tests..."
	@go test -v ./tjson/...

# Run XML tests
test-xml:
	@echo "Running XML tests..."
	@go test -v ./txml/...

# Run YAML tests
test-yaml:
	@echo "Running YAML tests..."
	@go test -v ./tyaml/...

# Run Node tests
test-node:
	@echo "Running Node tests..."
	@go test -v ./node/...

# Run Benchmark tests
test-bench:
	@echo "Running Benchmark tests..."
	@go test -bench=. -benchmem ./benchmark/...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Run specific package tests with coverage
test-json-cover:
	@echo "Running JSON tests with coverage..."
	@go test -cover ./tjson/...

test-xml-cover:
	@echo "Running XML tests with coverage..."
	@go test -cover ./txml/...

test-yaml-cover:
	@echo "Running YAML tests with coverage..."
	@go test -cover ./tyaml/...

test-node-cover:
	@echo "Running Node tests with coverage..."
	@go test -cover ./node/...

# Run all tests with verbose output
test-verbose:
	@echo "Running all tests with verbose output..."
	@go test -v ./...

# Run tests and generate coverage report
test-coverage-report:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out