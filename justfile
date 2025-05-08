# List of technologies for .gitignore generation (comma-separated)
technologies := "localstack,terraform,go,cursor"
ignore_entries := ".cursor,data,mobs"
# To add a new technology, append it to the technologies variable above (comma-separated)

# List all available commands
default:
    @just --list

# Create a new ADR
# Usage: just adr "Use bucket-per-tenant isolation model"
adr TITLE:
    @adr new "{{TITLE}}"

# Use standardized commit format
# Usage: just commit feat api "add tenant creation endpoint"
commit TYPE SCOPE MESSAGE:
    @scripts/commit.sh -t "{{TYPE}}" -s "{{SCOPE}}" -m "{{MESSAGE}}"

# Verify required tools (Docker, Docker Compose, LocalStack, Terraform, Go, AWS CLI)
deps:
    @echo "🔗 Checking required dependencies..."
    @scripts/checkdeps.sh

# Regenerate .gitignore for all technologies in the list
regenerate-gitignore:
    @echo "Regenerating .gitignore for: {{technologies}} and {{ignore_entries}}"
    @scripts/regenerate_gitignore.sh {{technologies}} {{ignore_entries}}
    @echo ".gitignore updated."

# Remove build artifacts and temporary files
clean:
    @echo "Cleaning build artifacts and temporary files..."
    rm -f mobs

# Compile Go application
build:
    @echo "Building Go application..."
    go build -o mobs ./src/cli

# Clean and build
rebuild: clean build

# Start LocalStack and deploy infrastructure
start:
    @echo "Starting LocalStack and deploying infrastructure..."
    @echo "Not yet implemented"

# Run all tests
test:
    @echo "Running all tests..."
    @go test -v ./test/domain/...

# Run static analysis
lint:
    @echo "Running linters..."
    @echo "Not yet implemented"

# Automatically fix linting issues
lint-correct:
    @echo "Fixing linting issues..."
    @echo "Not yet implemented"

# Format code according to standards
format:
    @echo "Formatting code..."
    @echo "Not yet implemented"

# Show formatting changes without applying
format-dry:
    @echo "Showing formatting changes (dry run)..."
    @echo "Not yet implemented"

# Run all checks (lint + tests)
check: lint test

# Run all automatic fixes (lint-correct + format)
fix: lint-correct format
