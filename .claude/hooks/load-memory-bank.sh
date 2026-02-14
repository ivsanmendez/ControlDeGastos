#!/bin/bash
set -euo pipefail

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-.}"
MEMORY_DIR="$PROJECT_DIR/memory-bank"

if [ ! -d "$MEMORY_DIR" ]; then
  echo "Memory bank directory not found at $MEMORY_DIR"
  exit 0
fi

echo "=== MEMORY BANK LOADED ==="
echo ""

# Load core documents in numbered order
for file in "$MEMORY_DIR"/[0-9][0-9]_*.md; do
  [ -f "$file" ] || continue
  echo "--- $(basename "$file") ---"
  cat "$file"
  echo ""
done

# Load subdirectory READMEs for awareness
for readme in "$MEMORY_DIR"/*/README.md; do
  [ -f "$readme" ] || continue
  subdir=$(basename "$(dirname "$readme")")
  echo "--- $subdir/README.md ---"
  cat "$readme"
  echo ""
done

echo "=== END MEMORY BANK ==="