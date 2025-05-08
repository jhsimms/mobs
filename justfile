# List of technologies for .gitignore generation (comma-separated)
technologies := "localstack,terraform,go,cursor"
ignore_entries := ".cursor"
# To add a new technology, append it to the technologies variable above (comma-separated)

# List all available commands
default:
    @just --list

# Use standardized commit format
# Usage: just commit feat api "add tenant creation endpoint"
commit TYPE SCOPE MESSAGE:
    @scripts/commit.sh "{{TYPE}}" "{{SCOPE}}" "{{MESSAGE}}"

# Verify required tools (Docker, Docker Compose, LocalStack, Terraform, Go, AWS CLI)
deps:
    @echo "🔗 Checking required dependencies..."
    @scripts/checkdeps.sh

# Regenerate .gitignore for all technologies in the list
regenerate-gitignore:
    @echo "Regenerating .gitignore for: {{technologies}}"
    @scripts/regenerate_gitignore.sh {{technologies}}
    @echo ".gitignore updated."

# Remove build artifacts and temporary files
clean:
    @echo "Cleaning build artifacts and temporary files..."
    @echo "Not yet implemented"

# Compile Go application
build:
    @echo "Building Go application..."
    @echo "Not yet implemented"

# Clean and build
rebuild: clean build

# Start LocalStack and deploy infrastructure
start:
    @echo "Starting LocalStack and deploying infrastructure..."
    @echo "Not yet implemented"

# Run tests
test:
    @echo "Running tests..."
    @echo "Not yet implemented"

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
