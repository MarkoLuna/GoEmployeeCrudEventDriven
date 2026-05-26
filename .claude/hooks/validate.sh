#!/usr/bin/env bash
# .claude/hooks/validate.sh
#
# Hook: run after the agent finishes a turn.
# Purpose: verify that both Go microservices still compile and all unit tests pass.
# Exit code != 0 will surface a warning in Claude Code's output.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SERVICES=("employee-service" "employee-consumer")

PASS=0
FAIL=0
ERRORS=()

run_step() {
  local label="$1"
  local dir="$2"
  local cmd="$3"

  echo "▶ [$label] $cmd  (in $dir)"
  if (cd "$dir" && eval "$cmd" 2>&1); then
    echo "  ✅ $label passed"
    ((PASS++)) || true
  else
    echo "  ❌ $label FAILED"
    ((FAIL++)) || true
    ERRORS+=("$label")
  fi
}

echo ""
echo "═══════════════════════════════════════════"
echo "  GoEmployeeCrudEventDriven — Validate Hook"
echo "═══════════════════════════════════════════"
echo ""

for SVC in "${SERVICES[@]}"; do
  DIR="$REPO_ROOT/$SVC"

  # 1. Compilation check
  run_step "$SVC: go vet" "$DIR" "go vet ./..."

  # 2. Build check (compile only, no binary output)
  run_step "$SVC: go build" "$DIR" "go build ./..."

  # 3. Unit tests (30-second timeout, no external dependencies needed)
  run_step "$SVC: go test" "$DIR" "go test -timeout 30s ./..."
done

echo ""
echo "───────────────────────────────────────────"
if [ "$FAIL" -eq 0 ]; then
  echo "  ✅ All $PASS checks passed."
else
  echo "  ❌ $FAIL check(s) failed: ${ERRORS[*]}"
  echo "     Fix the issues above before committing."
fi
echo "───────────────────────────────────────────"
echo ""

exit "$FAIL"
