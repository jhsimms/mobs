#!/usr/bin/env bash
set -euo pipefail

MAX_BODY_LENGTH=250

# Print usage information and exit
usage() {
    cat <<EOF
Usage: $0 -t TYPE -s SCOPE -m MESSAGE [files...]
  -t TYPE     Commit type (see below)
  -s SCOPE    Commit scope (e.g., component or file)
  -m MESSAGE  Commit message (max $MAX_BODY_LENGTH chars, no leading cap, no trailing period)
  files       Files to stage (optional, default: all changes)

Valid commit types:
  feat      A new feature
  fix       A bug fix
  docs      Documentation only changes
  style     Code style changes (formatting, missing semi colons, etc)
  refactor  Code changes that neither fix a bug nor add a feature
  test      Adding or correcting tests
  chore     Maintenance tasks (build, deps, etc)
EOF
    exit 1
}

# Print error message and exit
fail() {
    echo "Error: $1" >&2
    exit 1
}

# Validate commit type against allowed list
validate_commit_type() {
    local type="$1"
    local -a valid_types=("feat" "fix" "docs" "style" "refactor" "test" "chore")
    for t in "${valid_types[@]}"; do
        if [[ "$t" == "$type" ]]; then
            return 0
        fi
    done
    fail "Invalid commit type '$type'. Valid types: ${valid_types[*]}"
}

# Validate commit message length
validate_commit_message() {
    local message="$1"
    if (( ${#message} > MAX_BODY_LENGTH )); then
        fail "Commit message too long (${#message} chars, max $MAX_BODY_LENGTH)"
    fi
}

# Parse command-line arguments
parse_args() {
    local OPTIND opt
    TYPE=""
    SCOPE=""
    MESSAGE=""
    FILES=()
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
    if (( $# > 0 )); then
        FILES=("$@")
    fi
}

# Stage files for commit
stage_files() {
    git add .
}

# Prompt user for confirmation before committing
confirm_commit() {
    read -p "Continue with commit? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Commit aborted."
        exit 0
    fi
}

# Main entry point
main() {
    if [[ $# -eq 0 ]]; then
        usage
    fi

    parse_args "$@"

    if [[ -z "$TYPE" || -z "$SCOPE" || -z "$MESSAGE" ]]; then
        fail "TYPE, SCOPE, and MESSAGE are required."
    fi

    validate_commit_type "$TYPE"
    validate_commit_message "$MESSAGE"

    local commit_msg="$TYPE($SCOPE): $MESSAGE"

    stage_files

    echo "Committing with message: $commit_msg"
    echo "Changes to be committed:"
    git status --short

    confirm_commit

    git commit -m "$commit_msg"
    echo "Commit successful!"
}

main "$@"