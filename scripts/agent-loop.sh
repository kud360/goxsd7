#!/usr/bin/env bash
# agent-loop.sh — run one or more opencode trigger commands in sequence.
#
# Usage:
#   ./scripts/agent-loop.sh                    # default: develop
#   ./scripts/agent-loop.sh develop            # one develop iteration
#   ./scripts/agent-loop.sh plan,develop,retro # a trigger sequence
#   GOXSD_SEQUENCE=plan,develop ./scripts/agent-loop.sh
#
# Valid triggers: develop | ratchet | plan | retro
# (each maps to .opencode/command/<name>.md, run by the foreman agent)
#
# Designed for cron/launchd (see scripts/schedule/). A lock file prevents
# overlapping runs; logs go to .agent/logs/.
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SEQUENCE="${1:-${GOXSD_SEQUENCE:-develop}}"
LOG_DIR="$REPO_DIR/.agent/logs"
LOCK_FILE="$REPO_DIR/.agent/loop.lock"
OPENCODE_BIN="${OPENCODE_BIN:-opencode}"

mkdir -p "$LOG_DIR"

# Never overlap a running session. Wait (bounded) for the lock rather than
# skipping outright, so infrequent triggers (plan/retro) aren't starved by
# an aggressive develop cadence. Stale locks > 6h are cleared.
# GOXSD_LOCK_WAIT: max seconds to wait (default 900; 0 = skip immediately).
LOCK_WAIT="${GOXSD_LOCK_WAIT:-900}"
waited=0
while [ -e "$LOCK_FILE" ]; do
    if [ -n "$(find "$LOCK_FILE" -mmin +360 2>/dev/null)" ]; then
        echo "clearing stale lock ($LOCK_FILE)" >&2
        rm -f "$LOCK_FILE"
        break
    fi
    if [ "$waited" -ge "$LOCK_WAIT" ]; then
        echo "lock still held after ${waited}s ($LOCK_FILE); skipping" >&2
        exit 0
    fi
    sleep 15
    waited=$((waited + 15))
done
echo "$$" > "$LOCK_FILE"
trap 'rm -f "$LOCK_FILE"' EXIT

run_trigger() {
    local trigger="$1"
    case "$trigger" in
        develop|ratchet|plan|retro) ;;
        *) echo "unknown trigger: $trigger (want develop|ratchet|plan|retro)" >&2; return 1 ;;
    esac

    local ts log
    ts="$(date +%Y%m%d-%H%M%S)"
    log="$LOG_DIR/$ts-$trigger.log"
    echo "[$ts] trigger=$trigger log=$log"

    # `opencode run --command <name>` on recent versions; fall back to
    # sending the slash command as the message for older ones.
    if "$OPENCODE_BIN" run --help 2>&1 | grep -q -- '--command'; then
        (cd "$REPO_DIR" && "$OPENCODE_BIN" run --command "$trigger" --agent foreman) >"$log" 2>&1
        return
    fi
    (cd "$REPO_DIR" && "$OPENCODE_BIN" run --agent foreman "/$trigger") >"$log" 2>&1
}

IFS=',' read -ra TRIGGERS <<< "$SEQUENCE"
for t in "${TRIGGERS[@]}"; do
    run_trigger "$(echo "$t" | tr -d '[:space:]')"
done

# Keep the last 200 logs.
ls -1t "$LOG_DIR" | tail -n +201 | while read -r old; do rm -f "$LOG_DIR/$old"; done
