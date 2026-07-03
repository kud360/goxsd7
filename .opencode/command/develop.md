---
description: Run one development iteration (pick issue → oracle → mason → arbiter → chronicler)
agent: foreman
---

Run exactly one iteration of the develop loop from docs/WORKFLOW.md.

1. Ensure a clean, up-to-date tree (git status; git pull --rebase;
   git submodule update --init). A dirty tree gets logged and reverted first.
2. Pick the top unblocked `ready` issue: `gh issue list --label ready`.
   If none exists, delegate to @cartographer to replenish the ready set,
   log via @chronicler, and stop.
3. Get spec grounding from @oracle for every rule ID in the issue's scope.
4. Have @mason implement the smallest change that closes the issue.
   If public API changed, have @warden review before judging.
5. Have @arbiter judge (verdict + ratchet). Accept → commit using the
   format in AGENTS.md and close the issue. Reject twice → revert, comment
   findings on the issue, relabel `needs-replan`.
6. Have @chronicler append the session log entry. Push.

Stop after one issue. $ARGUMENTS
