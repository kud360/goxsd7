# goxsd7 — Agent Instructions

You are working on **goxsd7**, a conformance-grade XSD 1.1 processor in Go
(module `github.com/kud360/goxsd7`). This repo is developed primarily by AI
agents. Follow this file exactly; it wins over your own preferences.

## The one rule that outranks everything

**Never regress the ratchet.** Conformance expectations live in
`conformance/testdata/expectations/`. Scores only move up. If your change
makes a previously passing case fail, either fix it or revert your change.
Never edit an expectations file downward to make CI green.

## Ground truth

- The specs are LOCAL, in `docs/specs/md/`. Do not guess spec behavior from
  memory — grep the spec. Cite rule IDs (e.g. `cvc-complex-type.2.1`,
  `cos-st-restricts`) in code and commit messages when implementing them.
  - `xmlschema11-1.md` — Structures
  - `xmlschema11-2.md` — Datatypes (Appendix E hfn function definitions
    are the source of truth for builtin types)
  - `xpath20.md` — XPath 2.0
  - `xsd-precisionDecimal.md` — precisionDecimal
- Architecture: `docs/ARCHITECTURE.md`. Style: `docs/STYLE.md`
  (non-negotiable). Workflow: `docs/WORKFLOW.md`. Past-life lessons:
  `docs/LESSONS.md`. Roadmap: `docs/PLAN.md`.

## Commands

```sh
make build          # go build ./...
make test           # go test ./...
make vet            # go vet + golangci-lint (STYLE gate, incl. gofmt)
make conformance    # run W3C suite against current code
make ratchet        # conformance + update expectations upward only
make specs          # regenerate docs/specs/md from docs/specs/html
```

All four of `make build test vet conformance` must pass before any commit.

## Style (full rules in docs/STYLE.md — these are the headlines)

1. Happy path stays left; return early; **no `else` blocks**.
2. Every error is checked, wrapped with context, and mapped to a spec
   validation rule via `xsderr`. Errors carry schema file:line:column.
   No dropped errors inside loops — collect or return.
3. Deterministic output always. Never range over a map to produce output.
   Collections that reach users are slices in document order.
4. No state that can be derived from other state. No caches without a
   measured hot path.
5. Make illegal states unrepresentable: unexported fields + constructors,
   phased construction so cycle checks are unnecessary.
6. Comparison and facet capabilities are interfaces, not type switches.
7. Fail-open for unsupported XPath constructs (never false-reject), and
   track every fail-open site with a `// GAP(xpath): ...` comment.

## Workflow (full loop in docs/WORKFLOW.md)

Work is planned as GitHub issues (`gh issue list --label ready`). One issue
= one focused change = one commit. Commit message format:

```
<area>: <what changed> (#<issue>)

Spec: <rule ids touched>
Ratchet: <before> -> <after>
```

Append a dated entry to `docs/LOG/<year>-<month>.md` at the end of every
session (what was attempted, what worked, what surprised you).

## Personas

Specialized subagents are defined in `.opencode/agent/`:
**mason** (implements), **arbiter** (judges & runs the ratchet),
**oracle** (spec exegesis), **warden** (API/type-safety review),
**cartographer** (long-horizon planning, GitHub issues),
**chronicler** (logs & process retrospectives). Delegate to them per
`docs/WORKFLOW.md`; do not blur their roles.
