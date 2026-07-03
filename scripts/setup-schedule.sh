#!/usr/bin/env bash
# setup-schedule.sh — install/remove the scheduled agent loop for THIS
# clone, wherever it lives. Nothing about your directory layout needs to
# be committed: the launchd plists are rendered from
# scripts/schedule/goxsd7.plist.template at install time.
#
# Usage:
#   ./scripts/setup-schedule.sh check       # verify prerequisites only
#   ./scripts/setup-schedule.sh preset <steady|aggressive|relentless>
#       steady:     develop 2h; plan,develop daily 07:00; ratchet 22:00; retro sun 09:00
#       aggressive: develop 30m; plan 6h; ratchet daily 22:00; retro daily 09:00
#       relentless: develop 10m (back-to-back); plan 4h; ratchet 6h; retro daily 09:00
#       (replaces any previously installed goxsd7 agents)
#   ./scripts/setup-schedule.sh install <sequence> --every <seconds>
#   ./scripts/setup-schedule.sh install <sequence> --at <HH:MM> [--weekday 0-6]
#       (install runs `check` first; --skip-checks to bypass)
#   ./scripts/setup-schedule.sh uninstall [<sequence>|all]
#   ./scripts/setup-schedule.sh status
#   ./scripts/setup-schedule.sh cron [preset] # print crontab lines instead (any OS)
#
# <sequence> is a trigger or comma-list: develop | ratchet | plan | retro
# (see scripts/agent-loop.sh). Examples:
#   ./scripts/setup-schedule.sh install develop --every 7200
#   ./scripts/setup-schedule.sh install plan,develop --at 07:00
#   ./scripts/setup-schedule.sh install retro --at 09:00 --weekday 0
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMPLATE="$REPO_DIR/scripts/schedule/goxsd7.plist.template"
AGENTS_DIR="$HOME/Library/LaunchAgents"

die() { echo "setup-schedule: $*" >&2; exit 1; }

label_for() { echo "com.goxsd7.$(echo "$1" | tr ',' '-')"; }

check_sequence() {
    IFS=',' read -ra parts <<< "$1"
    for p in "${parts[@]}"; do
        case "$p" in
            develop|ratchet|plan|retro) ;;
            *) die "unknown trigger '$p' (want develop|ratchet|plan|retro)" ;;
        esac
    done
}

require_macos() {
    [ "$(uname)" = "Darwin" ] || die "launchd is macOS-only; use '$0 cron' instead"
}

# cmd_check verifies everything an unattended run needs. Missing tools and
# a missing test-suite submodule are fatal; things that only degrade the
# run (ollama down right now, gh not authed yet) are warnings.
cmd_check() {
    local fail=0 warn=0
    ok()      { echo "  ok    $*"; }
    bad()     { echo "  FAIL  $*"; fail=1; }
    warn_()   { echo "  warn  $*"; warn=1; }

    echo "checking prerequisites for $REPO_DIR"

    local tool path
    for tool in git go gh opencode golangci-lint ollama curl; do
        if ! path="$(command -v "$tool" 2>/dev/null)"; then
            case "$tool" in
                ollama|curl) warn_ "$tool not on PATH" ;;
                *)           bad  "$tool not on PATH" ;;
            esac
            continue
        fi
        ok "$tool ($path)"
        # launchd installs bake resolved tool dirs into the plist PATH;
        # cron does not — warn so crontab users add a PATH line.
        case "$path" in
            /opt/homebrew/bin/*|/usr/local/bin/*|/usr/bin/*|/bin/*) ;;
            *) warn_ "$tool lives outside the standard PATH — launchd installs handle this automatically; for cron add a PATH line to your crontab" ;;
        esac
    done

    if [ -f "$REPO_DIR/testdata/xsdtests/suite.xml" ]; then
        ok "W3C suite submodule initialized"
    else
        bad "W3C suite missing — run: git -C $REPO_DIR submodule update --init"
    fi

    if command -v gh >/dev/null 2>&1; then
        if gh auth status >/dev/null 2>&1; then
            ok "gh authenticated"
        else
            warn_ "gh not authenticated — run: gh auth login (agents need issue read/write + push)"
        fi
    fi

    if command -v curl >/dev/null 2>&1; then
        if curl -sf --max-time 3 http://localhost:11434/api/tags >/dev/null 2>&1; then
            ok "ollama serving on :11434"
            # The exact model opencode will ask for, e.g. "ollama/gemma4:31b-mlx".
            local model
            model="$(sed -n 's|.*"model":[[:space:]]*"ollama/\([^"]*\)".*|\1|p' "$REPO_DIR/opencode.json" | head -1)"
            if [ -z "$model" ]; then
                warn_ "could not read an ollama/* model from opencode.json — verify its \"model\" entry"
            elif command -v ollama >/dev/null 2>&1; then
                if ollama list 2>/dev/null | awk 'NR>1 {print $1}' | grep -qx "$model"; then
                    ok "model $model pulled"
                else
                    warn_ "model $model not in ollama list — run: ollama pull $model"
                fi
            fi
        else
            warn_ "ollama not responding on localhost:11434 — scheduled runs will no-op until it serves"
        fi
    fi

    if git -C "$REPO_DIR" push --dry-run >/dev/null 2>&1; then
        ok "git push works non-interactively"
    else
        warn_ "git push --dry-run failed — set up credentials that work unattended (SSH key / credential helper)"
    fi

    [ "$fail" -eq 0 ] || die "prerequisites missing (see FAIL lines above)"
    [ "$warn" -eq 0 ] || echo "warnings above won't block install, but fix them before trusting the schedule"
    echo "prerequisites look good"
}

# launchd_path builds the PATH baked into the rendered plist: the base
# system dirs plus wherever the required tools actually resolve on this
# machine (launchd does not inherit your shell's PATH).
launchd_path() {
    local path="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin" tool dir
    for tool in git go gh opencode golangci-lint ollama; do
        dir="$(command -v "$tool" 2>/dev/null | xargs -I{} dirname {} 2>/dev/null)" || true
        [ -n "$dir" ] || continue
        case ":$path:" in
            *":$dir:"*) ;;
            *) path="$dir:$path" ;;
        esac
    done
    echo "$path"
}

render() { # $1=label $2=sequence $3=schedule-xml
    sed -e "s|@@LABEL@@|$1|g" \
        -e "s|@@REPO@@|$REPO_DIR|g" \
        -e "s|@@SEQUENCE@@|$2|g" \
        -e "s|@@SCHEDULE@@|$3|g" \
        -e "s|@@PATH@@|$(launchd_path)|g" \
        "$TEMPLATE"
}

cmd_install() {
    require_macos
    local sequence="${1:-}"; shift || true
    [ -n "$sequence" ] || die "install needs a sequence (e.g. develop)"
    check_sequence "$sequence"

    local every="" at="" weekday="" skip_checks=0
    while [ $# -gt 0 ]; do
        case "$1" in
            --every)       every="${2:-}"; shift 2 ;;
            --at)          at="${2:-}"; shift 2 ;;
            --weekday)     weekday="${2:-}"; shift 2 ;;
            --skip-checks) skip_checks=1; shift ;;
            *) die "unknown flag $1" ;;
        esac
    done

    if [ "$skip_checks" -eq 0 ]; then
        cmd_check
        echo
    fi

    local schedule
    if [ -n "$every" ]; then
        [ -z "$at" ] || die "--every and --at are mutually exclusive"
        [[ "$every" =~ ^[0-9]+$ ]] || die "--every wants seconds, got '$every'"
        schedule="    <key>StartInterval</key>\\
    <integer>$every</integer>"
    elif [ -n "$at" ]; then
        [[ "$at" =~ ^([0-9]{1,2}):([0-9]{2})$ ]] || die "--at wants HH:MM, got '$at'"
        local hour="${BASH_REMATCH[1]}" minute="${BASH_REMATCH[2]}"
        schedule="    <key>StartCalendarInterval</key>\\
    <dict>\\
        <key>Hour</key><integer>$((10#$hour))</integer>\\
        <key>Minute</key><integer>$((10#$minute))</integer>"
        if [ -n "$weekday" ]; then
            schedule="$schedule\\
        <key>Weekday</key><integer>$weekday</integer>"
        fi
        schedule="$schedule\\
    </dict>"
    else
        die "install needs --every <seconds> or --at <HH:MM>"
    fi

    local label plist
    label="$(label_for "$sequence")"
    plist="$AGENTS_DIR/$label.plist"
    mkdir -p "$AGENTS_DIR" "$REPO_DIR/.agent"

    launchctl bootout "gui/$(id -u)/$label" 2>/dev/null || true
    render "$label" "$sequence" "$schedule" > "$plist"
    launchctl bootstrap "gui/$(id -u)" "$plist"
    echo "installed $label -> $plist"
    echo "logs: $REPO_DIR/.agent/  (launchd.out.log, launchd.err.log, logs/)"
}

cmd_preset() {
    require_macos
    local name="${1:-}"
    [ -n "$name" ] || die "preset needs a name (steady|aggressive|relentless)"
    cmd_check
    echo
    cmd_uninstall all >/dev/null
    case "$name" in
        steady)
            cmd_install develop --every 7200 --skip-checks
            cmd_install plan,develop --at 07:00 --skip-checks
            cmd_install ratchet --at 22:00 --skip-checks
            cmd_install retro --at 09:00 --weekday 0 --skip-checks
            ;;
        aggressive)
            cmd_install develop --every 1800 --skip-checks
            cmd_install plan --every 21600 --skip-checks
            cmd_install ratchet --at 22:00 --skip-checks
            cmd_install retro --at 09:00 --skip-checks
            ;;
        relentless)
            cmd_install develop --every 600 --skip-checks
            cmd_install plan --every 14400 --skip-checks
            cmd_install ratchet --every 21600 --skip-checks
            cmd_install retro --at 09:00 --skip-checks
            ;;
        *) die "unknown preset '$name' (steady|aggressive|relentless)" ;;
    esac
    echo
    cmd_status
}

cmd_uninstall() {
    require_macos
    local target="${1:-all}"
    local removed=0
    for plist in "$AGENTS_DIR"/com.goxsd7.*.plist; do
        [ -e "$plist" ] || continue
        local label
        label="$(basename "$plist" .plist)"
        if [ "$target" != "all" ] && [ "$label" != "$(label_for "$target")" ]; then
            continue
        fi
        launchctl bootout "gui/$(id -u)/$label" 2>/dev/null || true
        rm -f "$plist"
        echo "removed $label"
        removed=1
    done
    [ "$removed" -eq 1 ] || echo "nothing to remove"
}

cmd_status() {
    require_macos
    local found=0
    for plist in "$AGENTS_DIR"/com.goxsd7.*.plist; do
        [ -e "$plist" ] || continue
        found=1
        local label
        label="$(basename "$plist" .plist)"
        if launchctl print "gui/$(id -u)/$label" >/dev/null 2>&1; then
            echo "$label: loaded ($plist)"
        else
            echo "$label: NOT loaded ($plist)"
        fi
    done
    [ "$found" -eq 1 ] || echo "no goxsd7 agents installed"
    if [ -d "$REPO_DIR/.agent/logs" ]; then
        ls -1t "$REPO_DIR/.agent/logs" | head -3 | sed 's/^/recent run: /'
    fi
}

cmd_cron() {
    local preset="${1:-steady}"
    echo "# goxsd7 agent loop ($preset) — add with: crontab -e"
    case "$preset" in
        steady) cat <<EOF
0 */2 * * *   $REPO_DIR/scripts/agent-loop.sh develop
0 7  * * *    $REPO_DIR/scripts/agent-loop.sh plan,develop
0 22 * * *    $REPO_DIR/scripts/agent-loop.sh ratchet
0 9  * * 0    $REPO_DIR/scripts/agent-loop.sh retro
EOF
        ;;
        aggressive) cat <<EOF
*/30 * * * *  $REPO_DIR/scripts/agent-loop.sh develop
15 */6 * * *  $REPO_DIR/scripts/agent-loop.sh plan
0 22 * * *    $REPO_DIR/scripts/agent-loop.sh ratchet
0 9  * * *    $REPO_DIR/scripts/agent-loop.sh retro
EOF
        ;;
        relentless) cat <<EOF
*/10 * * * *  $REPO_DIR/scripts/agent-loop.sh develop
15 */4 * * *  $REPO_DIR/scripts/agent-loop.sh plan
30 */6 * * *  $REPO_DIR/scripts/agent-loop.sh ratchet
0 9  * * *    $REPO_DIR/scripts/agent-loop.sh retro
EOF
        ;;
        *) die "unknown preset '$preset' (steady|aggressive|relentless)" ;;
    esac
    cat <<'EOF'
# cron PATH is minimal; ensure opencode/go/gh/git/golangci-lint resolve, e.g.:
# PATH=/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin
EOF
}

case "${1:-}" in
    check)     cmd_check ;;
    preset)    shift; cmd_preset "$@" ;;
    install)   shift; cmd_install "$@" ;;
    uninstall) shift; cmd_uninstall "$@" ;;
    status)    cmd_status ;;
    cron)      shift || true; cmd_cron "$@" ;;
    *) sed -n '2,27p' "$0" | sed 's/^# \{0,1\}//'; exit 2 ;;
esac
