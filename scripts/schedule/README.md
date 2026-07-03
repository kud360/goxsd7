# Scheduling the agent loop

`scripts/agent-loop.sh` runs one or more triggers
(`develop|ratchet|plan|retro`, comma-separated). Nothing machine-specific
is committed: `scripts/setup-schedule.sh` renders and installs the
schedule for *your* clone, wherever it lives.

## Quick start (macOS, launchd)

```sh
./scripts/setup-schedule.sh check                               # verify prerequisites
./scripts/setup-schedule.sh install develop --every 7200        # every 2h
./scripts/setup-schedule.sh install plan,develop --at 07:00     # daily 07:00
./scripts/setup-schedule.sh install ratchet --at 22:00          # daily 22:00
./scripts/setup-schedule.sh install retro --at 09:00 --weekday 0  # sunday
./scripts/setup-schedule.sh status
./scripts/setup-schedule.sh uninstall all
```

Plists are rendered from `goxsd7.plist.template` into
`~/Library/LaunchAgents/com.goxsd7.<sequence>.plist`; run output lands in
`.agent/` inside the repo (gitignored).

## cron (any OS)

```sh
./scripts/setup-schedule.sh cron   # prints crontab lines for this clone
crontab -e                         # paste them
```

## Suggested cadences

| Cadence | Sequence | Why |
|---|---|---|
| Every 2 h | `develop` | steady stone-laying |
| Daily 07:00 | `plan,develop` | replenish ready issues, then work |
| Daily 22:00 | `ratchet` | end-of-day conformance health check |
| Sunday 09:00 | `retro` | weekly process self-improvement |

The loop's lock file makes overlapping fires safe (they skip).

## Prerequisites for unattended runs

`./scripts/setup-schedule.sh check` verifies all of these (and `install`
runs it automatically; `--skip-checks` to bypass):

- `git`, `go`, `gh`, `opencode` on PATH. launchd doesn't inherit your
  shell's PATH, so `install` bakes the resolved directories of these
  tools into the rendered plist; cron users get a warning to add a PATH
  line to their crontab instead.
- `ollama serve` responding with the model from opencode.json pulled
  (warning only — it just needs to be up when the schedule fires).
- `gh auth status` logged in (issue read/write, push).
- Git push credentials that work non-interactively (SSH key/credential
  helper; checked via `git push --dry-run`).
- W3C suite submodule initialized (`git submodule update --init`).
