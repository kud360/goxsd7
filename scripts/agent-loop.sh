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

# Never overlap a running session. The lock file holds the owner's PID;
# a lock is stale only if its owner is gone (age tells you nothing — a
# healthy session can outlive any age threshold, and clearing a live
# session's lock lets two sessions fight over one working tree).
# Wait (bounded) for a live lock rather than skipping outright, so
# infrequent triggers (plan/retro) aren't starved by a dense develop
# cadence. GOXSD_LOCK_WAIT: max seconds to wait (default 900; 0 = skip).
LOCK_WAIT="${GOXSD_LOCK_WAIT:-900}"
waited=0
while [ -e "$LOCK_FILE" ]; do
    holder="$(cat "$LOCK_FILE" 2>/dev/null || true)"
    if [ -z "$holder" ] || ! kill -0 "$holder" 2>/dev/null; then
        echo "clearing stale lock (owner ${holder:-unknown} not running)" >&2
        rm -f "$LOCK_FILE"
        break
    fi
    if [ "$waited" -ge "$LOCK_WAIT" ]; then
        echo "lock held by live pid $holder after ${waited}s; skipping" >&2
        exit 0
    fi
    sleep 15
    waited=$((waited + 15))
done
echo "$$" > "$LOCK_FILE"
trap 'rm -f "$LOCK_FILE"' EXIT

# Load the model and pin it in memory (keep_alive -1) before the session,
# so scheduled runs never pay a cold start and the pin holds regardless of
# how the ollama server's OLLAMA_KEEP_ALIVE env is set. Non-fatal: if
# ollama isn't up, the session itself will surface that.
warm_model() {
    local model
    model="$(sed -n 's|.*"model":[[:space:]]*"ollama/\([^"]*\)".*|\1|p' "$REPO_DIR/opencode.json" | head -1)"
    if [ -z "$model" ]; then
        return 0
    fi
    curl -sf --max-time 300 http://localhost:11434/api/generate \
        -d "{\"model\":\"$model\",\"keep_alive\":-1}" >/dev/null \
        || echo "warm_model: could not pin $model (ollama down?)" >&2
}

# Hard wall-clock budget per trigger. A session that hasn't converged by
# now never converges (observed: 10h of mason/arbiter ping-pong on one
# file); killing it loses less than letting it run. macOS ships no
# timeout(1), so babysit the PID ourselves.
# GOXSD_SESSION_TIMEOUT: max seconds per trigger (default 7200).
SESSION_TIMEOUT="${GOXSD_SESSION_TIMEOUT:-7200}"

run_with_timeout() {
    local secs="$1"; shift
    "$@" &
    local pid=$!
    (
        w=0
        while [ "$w" -lt "$secs" ] && kill -0 "$pid" 2>/dev/null; do
            sleep 30
            w=$((w + 30))
        done
        if kill -0 "$pid" 2>/dev/null; then
            echo "session exceeded ${secs}s; killing pid $pid" >&2
            kill -TERM "$pid" 2>/dev/null
            sleep 30
            kill -KILL "$pid" 2>/dev/null
        fi
    ) &
    local watchdog=$!
    local status=0
    wait "$pid" || status=$?
    kill "$watchdog" 2>/dev/null || true
    return "$status"
}

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
        (cd "$REPO_DIR" && run_with_timeout "$SESSION_TIMEOUT" \
            "$OPENCODE_BIN" run --command "$trigger" --agent foreman) >"$log" 2>&1
        return
    fi
    (cd "$REPO_DIR" && run_with_timeout "$SESSION_TIMEOUT" \
        "$OPENCODE_BIN" run --agent foreman "/$trigger") >"$log" 2>&1
}

warm_model

IFS=',' read -ra TRIGGERS <<< "$SEQUENCE"
for t in "${TRIGGERS[@]}"; do
    # A failed or timed-out trigger must not abort the rest of the
    # sequence (set -e): log it and move on.
    run_trigger "$(echo "$t" | tr -d '[:space:]')" \
        || echo "trigger $t exited non-zero" >&2
done

# Keep the last 200 logs.
ls -1t "$LOG_DIR" | tail -n +201 | while read -r old; do rm -f "$LOG_DIR/$old"; done
