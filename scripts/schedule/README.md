# Scheduling the agent loop

`scripts/agent-loop.sh` runs one or more triggers
(`develop|ratchet|plan|retro`, comma-separated). Nothing machine-specific
is committed: `scripts/setup-schedule.sh` renders and installs the
schedule for *your* clone, wherever it lives.

## Quick start (macOS, launchd)

```sh
./scripts/setup-schedule.sh check               # verify prerequisites
./scripts/setup-schedule.sh preset aggressive   # install a whole cadence set
./scripts/setup-schedule.sh status
./scripts/setup-schedule.sh uninstall all
```

Plists are rendered from `goxsd7.plist.template` into
`~/Library/LaunchAgents/com.goxsd7.<sequence>.plist`; run output lands in
`.agent/` inside the repo (gitignored).

## Presets

| Preset | develop | plan | ratchet | retro | Character |
|---|---|---|---|---|---|
| `steady` | every 2 h | daily 07:00 (+develop) | daily 22:00 | Sunday 09:00 | background hum; machine mostly idle |
| `aggressive` | every 30 min | every 6 h | daily 22:00 | daily 09:00 | model busy most of the day; ~30–40 sessions/day |
| `relentless` | every 10 min | every 4 h | every 6 h | daily 09:00 | back-to-back sessions; GPU pinned continuously |

Individual installs remain available for custom mixes:

```sh
./scripts/setup-schedule.sh install develop --every 1800
./scripts/setup-schedule.sh install retro --at 09:00 --weekday 0
```

### How aggressive cadences behave

- One session at a time, always: the loop's lock file serializes runs.
  A fire that finds the lock held now **waits up to 15 minutes**
  (`GOXSD_LOCK_WAIT` seconds to tune, `0` = skip immediately), then
  skips — so infrequent triggers (plan/retro) aren't starved by a dense
  develop cadence, and a 30-min develop cadence with long sessions
  degrades gracefully into back-to-back execution.
- Faster cadences don't make sessions smarter, they just retry sooner
  after failures and idle less between issues. If sessions routinely run
  longer than the develop interval, the extra fires simply skip — cost
  is nil, but `relentless` keeps the GPU warm essentially 24/7 (heat,
  fan, power on a laptop).
- More sessions/day also means more commits/pushes and more issue churn;
  keep an eye on the `ready` queue depth — if the cartographer can't
  keep 5–10 issues ready, develop fires will spend themselves on
  planning instead (by design).

## cron (any OS)

```sh
./scripts/setup-schedule.sh cron aggressive   # prints crontab lines for this clone
crontab -e                                    # paste them
```

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
