#!/usr/bin/env bash
# .claude/commands/tidy.sh
# Run `go mod tidy` across all Go modules in the repo
#
# Usage:
#   /project:tidy

set -e

ROOT="$(git rev-parse --show-toplevel)"

for mod_dir in "$ROOT/common" "$ROOT/employee-service" "$ROOT/employee-consumer"; do
  echo "=== go mod tidy in $mod_dir ==="
  (cd "$mod_dir" && go mod tidy)
done

echo "All modules tidied."
