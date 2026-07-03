# Lessons from goxsd5

Distilled from the previous attempt
(https://github.com/kud360/goxsd5/blob/main/NOTES.md). These are paid-for
lessons; re-learning one is a process failure worth a chronicler entry.

## Architecture

1. **The core model must be a pure leaf.** goxsd5 started with `xsd`
   depending on datatype implementations and regex; unwinding that was a
   major refactor. goxsd7 starts leaf-first (see ARCHITECTURE.md).
2. **`Value` must be open (`any` + capability interfaces).** A sealed value
   interface blocks user backends. Capability detection through small
   interfaces beats type switches scattered across packages.
3. **Minimal capability interfaces at boundaries** (`ElementByName`-style
   views) are safer than exposing whole objects.
4. **Marker-interface infoset** scales to multiple formats (XML, JSON);
   sealed node types don't.
5. **Derivable state is a liability.** The memoized effective-facets cache
   was "hot path" superstition; on-demand merging was safe and simpler.

## Spec conformance traps (verified the hard way)

6. **Two regex flavors.** Pattern facets use XSD Part 2 grammar; assertion
   functions (`fn:matches` etc.) use F&O grammar. They differ on `^`/`$`
   (real anchors in F&O) and on `.` (F&O excludes only `\n`; pattern also
   excludes `\r`). One parser, two flavor flags — never one flavor for both.
7. **Union validation uses DirectMembers, not flattened members** — an
   intervening restriction's pattern/enum facets must be checked, and
   pattern matching normalizes with the *validating member's* whiteSpace.
8. **Assertions live at every level of the variety chain.** Lists: per item
   against the item type (including its assertions) before list-level
   assertions. Unions: the chosen member's assertions count too. Missing
   any level causes systematic false-accepts.
9. **Empty content is stricter than element-only content**: a complex type
   with `<xs:sequence/>` forbids even whitespace (cvc-complex-type.2.1).
10. **EDC with wildcards:** the *post-xsi:type governing* type must be
    validly derived from the locally declared type. Two plausible-sounding
    alternatives are wrong (goxsd5 tried both first).
11. **Content matching must be greedy** — consume all occurrence matches
    before exiting a loop, or invalid children leak to open-content
    wildcards. UPA determinism means no backtracking is needed. The
    remaining hard part: a wildcard must not absorb what the explicit model
    can still consume; plan for a real deterministic automaton.
12. **IDC namespace context is stateful.** `xpathDefaultNamespace` applies
    to element steps but not attribute steps; thread namespace context
    through the matcher chain from day one.
13. **`xs:override` needs explicit target tracking** — schema-level
    defaults, replacement ownership, and groups declared inside the
    override all attach to the *overridden* document. Uniform rules fail
    three different ways.
14. **Value constraints (default/fixed) interact with ID harvesting and
    empty content**; the element's *parent* is needed for cvc-id binding.
    Thread the parent through assessment from the start.
15. **XPath variables need types.** String-only `$value` bindings are
    insufficient; carry `{Lexical, Kind}` typed atoms.

## Strategy

16. **Fail-open is the right way to ship a partial XPath** (never
    false-reject), but only with systematic, greppable gap tracking and a
    ratchet closing them — otherwise false-accepts silently pile up.
17. **Ratchet expectations belong in version control.** Diffs make
    regressions obvious; blame makes them bisectable. Never re-baseline
    downward; never re-baseline without understanding the delta.
18. **Small fix → re-baseline → commit** beats batching. Independent
    reverts saved goxsd5 repeatedly.
19. **Throwaway diagnostic tests are a first-class tool.** Env-gated dumps
    of all failing cases (`GAPDIAG=1`) beat reading test output for
    hundreds of cases. Write them crude, delete them after.
20. **Fuzzing finds panics, not logic bugs.** Worth having for
    regex/value/XML parsing safety; don't expect conformance wins from it.
21. **Streaming from the start**: no `io.ReadAll` on documents; offset
    index for line/col instead of storing line starts.
22. **Some W3C cases exercise spec bugs** (goxsd5: `saxon all308`,
    `complex018`). Record as expected-divergence with justification rather
    than contorting the implementation; never let them cause false
    rejections elsewhere.

## Process (earned in goxsd7 — see docs/LOG/2026-07.md, 2026-07-03)

23. **Never destroy a dirty tree; rescue it.** `git clean -fd` at session
    start deleted 10 hours of uncommitted overnight work. Uncommitted
    changes are a crashed session's output: `git stash push -u -m
    "rescue <ts>"` and log it.
24. **The session log rides in the session commit.** A log entry written
    after the commit is left uncommitted, the next session sees a dirty
    tree it didn't cause, and the entry gets reverted — session history
    silently evaporates.
25. **Two rejections is the convergence horizon.** 11 mason attempts and
    7 arbiter judgments on one 428-line file converged no further than
    attempt 2 did on the retry the next day. After two rejects: rescue,
    comment, `needs-replan`, stop.
26. **Spec data tables are generated, never hand-typed.** An LLM
    hand-writing 49 datatypes' regexes token-by-token is the slowest,
    most error-prone path and can never satisfy byte-identical
    regeneration (D1). Write the generator; commit generator + output.
27. **Context size is a memory budget, not a convenience.** A 262144
    num_ctx on a 64 GiB machine hit 108 GiB peak (deep swap), tripled
    per-token cost, and made built-in MTP speculation a net loss so it
    disabled itself. Cap num_ctx to what RAM affords.
