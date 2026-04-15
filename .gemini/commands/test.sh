#!/usr/bin/env bash
# .gemini/commands/test.sh
# Run tests for all services (or a specific one)
#
# Usage (from Claude slash command):
#   /project:test              → runs tests in both services
#   /project:test service      → runs tests in employee-service only
#   /project:test consumer     → runs tests in employee-consumer only

set -e

ROOT="$(git rev-parse --show-toplevel)"

run_tests() {
  local dir="$1"
  echo "=== Running tests in $dir ==="
  (cd "$ROOT/$dir" && go test ./... -v -count=1)
}

case "${1:-all}" in
  service)   run_tests "employee-service" ;;
  consumer)  run_tests "employee-consumer" ;;
  *)
    run_tests "employee-service"
    run_tests "employee-consumer"
    ;;
esac
