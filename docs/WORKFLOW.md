# Development Workflow

goxsd7 is developed by scheduled agent sessions (opencode + Gemma 4 via
Ollama). Long-horizon memory lives in **GitHub issues** (plans) and
**docs/LOG/** (history). A session is short; the repo is the brain.

## The cast

| Agent | File | Role |
|---|---|---|
| **mason** | `.opencode/agent/mason.md` | Implements one issue at a time |
| **arbiter** | `.opencode/agent/arbiter.md` | Judges changes; owns the ratchet verdict |
| **oracle** | `.opencode/agent/oracle.md` | Spec exegesis; answers only from `docs/specs/md` |
| **warden** | `.opencode/agent/warden.md` | API/type-safety review; illegal states unrepresentable |
| **cartographer** | `.opencode/agent/cartographer.md` | Long-horizon planning; owns GitHub issues/milestones |
| **chronicler** | `.opencode/agent/chronicler.md` | Session logs; meta-process retrospectives |

Roles do not blur: mason never re-baselines the ratchet, arbiter never
implements fixes, oracle never writes code.

## The develop loop (`/develop`, the default scheduled trigger)

1. **Pick** — `gh issue list --label ready --limit 10`; take the highest
   priority issue whose dependencies are closed. No ready issue → run the
   cartographer instead and stop.
2. **Understand** — read the issue; ask the **oracle** for the exact spec
   clauses and rule IDs in scope. The oracle's citation goes in the commit.
3. **Implement** — **mason** makes the smallest change that closes the
   issue, per docs/STYLE.md. New/changed public API → **warden** reviews
   before proceeding.
4. **Judge** — **arbiter** runs `make build test vet conformance`,
   reviews the diff against STYLE.md, and issues a verdict:
   - *accept* → arbiter runs `make ratchet` (upward only).
   - *reject* → one repair round by mason, then re-judge. Second reject →
     stash the work (`git stash push -u -m "rescue <ts>"`), comment
     findings on the issue, relabel `needs-replan`. Two rejections is
     the hard cap — a third attempt has never converged.
5. **Record & commit** — **chronicler** appends to
   `docs/LOG/<year>-<month>.md` FIRST; then one commit carries the code
   and the log entry together; close or comment the issue
   (`gh issue close/comment`); push. The tree is clean after every push —
   a session that leaves docs/LOG uncommitted has failed.

Budget: one issue per session. Nothing works? A rescue stash + a good
issue comment is a successful session.

Never destroy a dirty tree: `git clean -fd` / `git restore .` deleted 10
hours of uncommitted session output once (2026-07-03). Uncommitted
changes at session start get stashed as `rescue <timestamp>` and logged,
so a human can triage them later.

## Other triggers

- **`/ratchet`** — arbiter only: run conformance, report movement, ratchet
  upward, investigate & file issues for any regression.
- **`/plan`** — cartographer: reconcile GitHub issues with reality (close
  stale, split oversized, order by dependency, keep 5–10 `ready`).
- **`/retro`** — chronicler: read the last ~2 weeks of LOG + git history;
  find recurring friction; propose concrete edits to WORKFLOW/STYLE/agent
  prompts as a PR. This is the self-improvement loop.

A typical schedule interleaves them, e.g. `plan,develop,develop,develop,retro`
(see `scripts/agent-loop.sh` and `scripts/schedule/`).

## GitHub conventions

- **Labels**: `ready` (unblocked, sized for one session), `blocked`,
  `needs-replan`, `epic`; areas `area/{model,parser,value,xpath,validate,codegen,codec,regex,loader,conformance,meta}`;
  kinds `kind/{feature,gap,bug,refactor,process}`.
- **Milestones** mirror docs/PLAN.md (M1, M2, …).
- Issue body must contain: goal, spec references (rule IDs), acceptance
  criteria (which conformance cases / tests prove it), and dependencies.
  If an agent can't start it from the body alone, the body is incomplete.
- `// GAP(...)` comments and fail-open sites get tracking issues
  (`kind/gap`) so nothing fails open silently forever.

## Commit format

```
<area>: <what changed> (#<issue>)

Spec: <rule ids, or "n/a">
Ratchet: <passed/total before> -> <after>   (or "unchanged")
```

Small, focused, independently revertible. Ratchet expectation updates are
part of the same commit as the fix that earned them.

## The ratchet (the heart of the process)

- `conformance/testdata/expectations/*.txt`: one line per W3C test case,
  `<case-id> <expected-outcome>`, sorted, committed.
- `make conformance` fails if any case does worse than its expectation.
- `make ratchet` rewrites expectations for cases that now do better,
  refuses to write anything worse.
- Every expectation change must be explainable; "it flipped and I don't
  know why" blocks the commit and becomes an issue.

## Debugging playbook (for agents)

- Reproduce one failing conformance case in isolation before touching code
  (the harness supports single-case runs; see conformance/README.md).
- Turn on scoped debug logs (`GOXSD_DEBUG=validate,xpath go test ...`) —
  messages carry rule IDs and locations by design.
- For bulk failure analysis, write an env-gated throwaway diagnostic test
  (`zz_diag_test.go`, gated on `DIAG=1`), harvest the pattern, delete it.
- Grep the spec (`docs/specs/md/`), not your memory. Quote clauses in
  issues and commits.
