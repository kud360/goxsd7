---
description: Run one development iteration (pick issue → oracle → mason → arbiter → chronicler)
agent: foreman
---

Run exactly one iteration of the develop loop from docs/WORKFLOW.md.

1. Ensure an up-to-date tree (git status; git pull --rebase;
   git submodule update --init). A dirty tree gets logged and RESCUED
   first: `git stash push -u -m "rescue <timestamp>"` — never
   `git clean`/`git restore`/`git checkout --`; uncommitted work is a
   crashed session's output, not garbage.
2. Pick the top unblocked `ready` issue: `gh issue list --label ready`.
   If none exists, delegate to @cartographer to replenish the ready set,
   log via @chronicler, and stop.
3. Get spec grounding from @oracle for every rule ID in the issue's scope.
4. Have @mason implement the smallest change that closes the issue.
   If public API changed, have @warden review before judging.
5. Have @arbiter judge (verdict + ratchet). At most TWO rejections total —
   count them. Second reject → stash the work as a rescue, comment
   findings on the issue, relabel `needs-replan`, skip to step 6.
6. Have @chronicler append the session log entry, THEN make one commit
   containing the code and the log entry together, using the format in
   AGENTS.md. Close the issue. Push. The tree must be clean after the
   push — docs/LOG is never left uncommitted.

Stop after one issue. $ARGUMENTS
