#!/usr/bin/env bash
# .gemini/commands/swagger.sh
# Regenerate Swagger docs for a service
#
# Usage:
#   /project:swagger service    → regenerates docs for employee-service
#   /project:swagger consumer   → regenerates docs for employee-consumer

set -e

ROOT="$(git rev-parse --show-toplevel)"

regen_swagger() {
  local dir="$1"
  echo "=== Regenerating Swagger docs in $dir ==="
  # Ensure go and swag are in PATH
  export PATH=$PATH:/home/marko/go/bin:/usr/local/go/bin
  (cd "$ROOT/$dir" && swag init --dir cmd/api-server,internal --parseDependency --parseInternal)
  echo "Done. Docs written to $dir/docs/"
}

case "${1:-}" in
  service)   regen_swagger "employee-service" ;;
  consumer)  regen_swagger "employee-consumer" ;;
  *)
    echo "Usage: swagger.sh [service|consumer]"
    exit 1
    ;;
esac
