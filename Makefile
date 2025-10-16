# Copyright (C) 2025 Michael Graff
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, version 3.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.

.PHONY: all check test build lint license-check fmt clean generate

# Default target: run tests and build
all: test build

# Build the binary
build:
	@echo "Building meshmgr..."
	@mkdir -p bin
	go build -o bin/meshmgr ./cmd/meshmgr

# Generate code (SQLC, etc.)
generate:
	@echo "Generating code..."
	go generate ./...

# Run tests
test: generate
	@echo "Running tests..."
	go test ./...

# Check everything: tests, lint, and license headers
check: test lint license-check
	@echo "All checks passed!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

# Check for license headers
license-check:
	@echo "Checking license headers..."
	@missing=$$(find . -name "*.go" -not -path "./vendor/*" -exec grep -L "GNU Affero General Public License" {} \;); \
	if [ -n "$$missing" ]; then \
		echo "Missing license headers in:"; \
		echo "$$missing"; \
		exit 1; \
	fi
	@echo "All files have license headers."

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean
