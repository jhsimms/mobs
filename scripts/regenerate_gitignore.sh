#!/usr/bin/env bash
set -euo pipefail
if [ $# -lt 1 ]; then
  echo "Usage: $0 <comma-separated-technologies> [comma-separated-additional-entries]"
  exit 1
fi

TECHNOLOGIES="$1"
ADDITIONAL_ENTRIES="$2"

curl -sL "https://www.toptal.com/developers/gitignore/api/${TECHNOLOGIES}" > .gitignore

if [ -n "$ADDITIONAL_ENTRIES" ]; then
  IFS=',' read -ra ENTRIES <<< "$ADDITIONAL_ENTRIES"
  echo "" >> .gitignore
  echo "# Additional entries" >> .gitignore
  for entry in "${ENTRIES[@]}"; do
    echo "$entry" >> .gitignore
  done
fi
