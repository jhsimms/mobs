#!/usr/bin/env bash
# scripts/checkdeps.sh
# Checks for required dependencies and prints all missing ones.

set -e
missing=0

deps=(docker docker-compose localstack terraform go aws)

for dep in "${deps[@]}"; do
    if ! command -v "$dep" > /dev/null; then
        echo "❌ Missing dependency: $dep"
        missing=1
    fi
done

if [ $missing -eq 0 ]; then
    echo "✅ All dependencies are installed!"
else
    echo "❌ Some dependencies are missing. Please install them and try again."
    exit 1
fi
