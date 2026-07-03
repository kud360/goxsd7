# Roadmap

Milestones are ordered by dependency; each becomes a GitHub milestone with
issues carved out by the cartographer. "Done" always means: builds, tested,
style-clean, ratchet not regressed, logged.

## M0 — Scaffold (done)

Repo layout, agent personas, opencode config, specs downloaded + converted,
W3C suite submodule, ratchet mechanics, package contracts (`doc.go`s),
`xsderr` and `loader` foundations.

## M1 — Spec infrastructure

- `tools/hfnextract`: extract Appendix E function definitions (hfn) and the
  builtin-datatype tables from `docs/specs/md/xmlschema11-2.md` into
  checked-in Go data (`builtin/gen_*.go`): per-type lexical grammar refs,
  facet applicability, ordered/bounded/cardinality/numeric properties,
  primitive ancestry. (#1, #2, #3)
- Acceptance: generated tables cover all 49 builtin types; regeneration is
  byte-identical (D1).

## M2 — XML foundation

- `parser/xmltree`: streaming reader, line/column on every node, bounded
  memory, namespace resolution, encoding handling. (#4, #5)
- `xsderr`: complete (already scaffolded) + rule-ID catalog generation from
  the specs (nice-to-have).
- Acceptance: fuzz target runs clean; positions verified by unit tests.

## M3 — Datatypes vertical slice (before structures!)

Facet pipeline + `builtin/strict` for a first tranche of types
(string family, decimal/integer family, boolean), driven by M1 tables.
Then the rest: float/double, date/time 7-property model, durations, binary,
anyURI, QName/NOTATION, precisionDecimal. (#6)
- Wire the **datatype-focused subset** of the W3C suite (msMeta/saxonMeta
  simple-type sets) into the ratchet → first real conformance numbers.
- Acceptance: ratchet lane live and climbing; user-defined restriction of a
  builtin works through the public API.


## M4 — Schema parsing & component model

- `parser`: schema documents → `xsd` components, phased construction
  (parse / resolve / finalize), `src-*` and `st/ct-props-correct` checks,
  include/import/redefine/override, chameleon includes.
- `loader` wired end to end (URL, location hints, catalogs).
- Acceptance: schema-validity ratchet lane live across the full suite.

## M5 — Instance validation core

- `validate`: infoset adapters, simple content via the pipeline, complex
  content with greedy deterministic matching, attributes, wildcards, EDC,
  ID/IDREF, IDC (with namespace-context threading — LESSONS 12), value
  constraints (LESSONS 14).
- Acceptance: instance-validity ratchet lane live.

## M6 — XPath subset → assertions & CTA

- `xpath`: CTA restricted subset first, then assertion essentials; typed
  `$value`; fail-open with `// GAP(xpath)` tracking; `regex.FO` flavor for
  F&O functions (LESSONS 6).
- Acceptance: assertion/CTA conformance cases climbing; gap inventory
  auto-generated into issues.

## M7 — Codegen

- `codegen`: deterministic Go type emission; strict/native backend choice;
  users can derive their own types from builtins in generated code.
- xs:choice → sealed interface + branch types (STYLE T2 closed-sum
  exception); required/optional/nillable shapes reviewed by the warden.
- Namer component: XSD → Go identifier wrangling; anonymous types named
  from schema ancestor context, deterministic document-order
  disambiguation; unit-tested against pathological schemas (same local
  names at many depths, Go keywords, case-only differences).
- Fix the emitter seam API (`value.Emitter`) and wire the first backend
  emitters (strict + native scalars).
- Acceptance: golden-file tests; generated code for the W3C suite schemas
  compiles; emitted decode functions carry type QName + schema Loc
  headers.

## M8 — Codec (dataset ser/de)

- `codec`: schema-directed decode into generated/reflective values,
  canonical encode; round-trip property tests. Runtime hot path uses
  ParseBytes/AppendCanonical appender conventions.
- Differential test rig: generated fast path vs runtime path on identical
  inputs must yield identical values and identical error rule IDs.
- Decode errors carry pipeline stage + type QName + input fragment +
  Loc/byte offset; `GOXSD_DEBUG=codec` stage tracing.
- Acceptance: round-trip = identity on canonical forms across suite
  instances; allocation benchmarks in CI (`testing.AllocsPerRun`) with
  zero-alloc scalar decode on the native fast path.

## M9 — `builtin/native` + backend polish

- Native Go backend; `value.Override` composition hardened; documented
  deviation table strict↔native; example custom user backend in
  `examples/`.

## Continuous (no milestone)

- Ratchet climbing; XPath gap closure; fuzz targets as parsers land;
  XPath 2.0 test suite lane when M6 matures; process retros.

## Explicit non-goals (for now)

- Mutation/editing APIs (`xsdedit`-style) — out of scope.
- XSD 1.0-only quirks compatibility mode.
- JSON infoset adapter (design for it, don't build it yet).
