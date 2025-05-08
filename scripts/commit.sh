#!/usr/bin/env bash
set -eu

MAX_BODY_LENGTH=250

usage() {
    echo "Usage: $0 -t TYPE -s SCOPE -m MESSAGE [files...]"
    echo "  -t TYPE     Commit type (feat, fix, docs, style, refactor, test, chore)"
    echo "  -s SCOPE    Commit scope (e.g., component or file)"
    echo "  -m MESSAGE  Commit message (max $MAX_BODY_LENGTH chars, no leading cap, no trailing period)"
    echo "  files       Files to stage (optional, default: all changes)"
    exit 1
}

# Default values
TYPE=""
SCOPE=""
MESSAGE=""
FILES=()

# Parse arguments
while getopts ":t:s:m:h" opt; do
  case $opt in
    t) TYPE="$OPTARG" ;;
    s) SCOPE="$OPTARG" ;;
    m) MESSAGE="$OPTARG" ;;
    h) usage ;;
    *) usage ;;
  esac
done
shift $((OPTIND -1))

# Remaining args are files
if [ $# -gt 0 ]; then
    FILES=("$@")
fi

# Validate required arguments
if [[ -z "$TYPE" || -z "$SCOPE" || -z "$MESSAGE" ]]; then
    echo "Error: TYPE, SCOPE, and MESSAGE are required."
    usage
fi

# Validate commit type
valid_types=("feat" "fix" "docs" "style" "refactor" "test" "chore")
valid_type=false
for t in "${valid_types[@]}"; do
    if [ "$t" = "$TYPE" ]; then
        valid_type=true
        break
    fi
done
if [ "$valid_type" = false ]; then
    echo "Error: Invalid commit type '$TYPE'"
    echo "Valid types: ${valid_types[*]}"
    exit 1
fi

# Validate commit message length
if [ ${#MESSAGE} -gt $MAX_BODY_LENGTH ]; then
    echo "Error: Commit message too long (${#MESSAGE} chars, max $MAX_BODY_LENGTH)"
    exit 1
fi

# Create the formatted commit message
commit_msg="$TYPE($SCOPE): $MESSAGE"

# Stage changes
if [ ${#FILES[@]} -eq 0 ]; then
    echo "No files specified. Staging all changes."
    git add .
else
    git add "${FILES[@]}"
fi

# Show what's being committed
echo "Committing with message: $commit_msg"
echo "Changes to be committed:"
git status --short

# Ask for confirmation
read -p "Continue with commit? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Commit aborted."
    exit 0
fi

# Commit with the formatted message
git commit -m "$commit_msg"

echo "Commit successful!"