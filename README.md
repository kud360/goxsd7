# goxsd7

A conformance-grade **XSD 1.1** processor for Go: schema parser, instance
validator, code generator, and dataset serializer/deserializer — measured
against the W3C XML Schema test suite and ratcheted upward, one commit at a
time.

This repository is **self-developing**: most work is performed by scheduled
local AI agents (opencode + Gemma 4) following the workflow in
[docs/WORKFLOW.md](docs/WORKFLOW.md), planned long-horizon via GitHub issues,
and continuously improved via retrospectives.

## Goals

- **Full XSD 1.1 conformance**, verified against the W3C test suite
  (`testdata/xsdtests`, git submodule) with a ratchet that only moves up.
- **Clean lexical/value-space separation** with a composable facet pipeline
  users can extend with their own types.
- **Builtin types bootstrapped from the spec** — generated from the function
  definitions (hfn) in XSD 1.1 Part 2 Appendix E, not hand-transcribed.
- **Two value backends**: `builtin/strict` (spec-exact, arbitrary precision,
  precisionDecimal) and `builtin/native` (Go-friendly: `int64`, `float64`,
  `time.Time`).
- **Errors that map to spec validation rules**, always decorated, always
  carrying schema file/line/column.
- **Deterministic output everywhere**; illegal states unrepresentable.

## Layout

| Path | Purpose |
|---|---|
| `xsd/` | Pure-leaf schema component model |
| `value/` | Value-space contracts + facet pipeline |
| `builtin/strict`, `builtin/native` | The two builtin type backends |
| `parser/`, `parser/xmltree` | Schema parsing; position-tracking XML reader |
| `loader/` | Schema resolution interfaces (URL, location) + helpers |
| `regex/` | XSD pattern & XPath F&O regex flavors |
| `xpath/` | XPath 2.0 engine (grown from the XSD subset, fail-open) |
| `validate/` | Instance validation |
| `codegen/` | Go code generation from schemas |
| `codec/` | Dataset (instance document) serializer/deserializer |
| `conformance/` | W3C suite harness + ratchet |
| `xsderr/` | Spec-rule-mapped, located, decorated errors |
| `docs/specs/` | Local specs (HTML + generated Markdown) |
| `tools/spec2md` | W3C spec HTML → Markdown converter |
| `.opencode/` | Agent personas and trigger commands |

## Getting started

```sh
git submodule update --init          # W3C test suite (~215 MB)
make specs                           # regenerate docs/specs/md from HTML
make build test                      # compile + unit tests
make conformance                     # run the ratchet (once implemented)
```

To run the agent loop manually:

```sh
./scripts/agent-loop.sh develop            # one development iteration
./scripts/agent-loop.sh plan,develop,retro # a custom trigger sequence
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md),
[docs/STYLE.md](docs/STYLE.md) (non-negotiable), and
[docs/LESSONS.md](docs/LESSONS.md) (distilled from the goxsd5 attempt).