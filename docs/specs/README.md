# Local specifications

Ground truth for all spec questions. **Grep these; never answer from
memory.** Regenerate Markdown with `make specs` (runs `tools/spec2md`).

| Markdown | Source | What |
|---|---|---|
| `md/xmlschema11-1.md` | [W3C REC xmlschema11-1](https://www.w3.org/TR/xmlschema11-1/) | XSD 1.1 Part 1: Structures |
| `md/xmlschema11-2.md` | [W3C REC xmlschema11-2](https://www.w3.org/TR/xmlschema11-2/) | XSD 1.1 Part 2: Datatypes |
| `md/xpath20.md` | [W3C REC xpath20](https://www.w3.org/TR/xpath20/) | XPath 2.0 (Second Edition) |
| `md/xsd-precisionDecimal.md` | [W3C WD xsd-precisionDecimal](https://www.w3.org/TR/xsd-precisionDecimal/) | precisionDecimal datatype |

The `html/` directory holds the pristine downloads (re-fetch with
`scripts/fetch-specs.sh`); `md/` is generated and committed so agents can
grep without tooling.

## How to find things (anchors survive conversion as `<a id="...">`)

- **Validation rules**: grep the rule ID directly — `cvc-complex-type`,
  `cos-st-restricts`, `src-resolve`, `derivation-ok-restriction`, …
- **The "hfn" function definitions** (normative lexical/canonical mappings
  that bootstrap our builtins): Datatypes **Appendix E**, anchor
  `id="ap-funcDefs"`. Individual functions are anchored `id="f-<name>"`
  (`f-decimalLexmap`, `f-booleanCanmap`, `f-dayTimeDurationMap`, …) and
  each builtin's mapping section is anchored
  `id="<type>-lexical-mapping"` / `id="<type>-canonical-mapping"`.
- **Facets**: `id="rf-facets"` region in Datatypes; per-facet anchors like
  `id="rf-pattern"`, `id="rf-minInclusive"`.
- **Builtin datatypes**: per-type sections anchored by type name
  (`id="decimal"`, `id="dateTime"`, …) in Datatypes §3.
- **XPath grammar**: `xpath20.md`, EBNF productions are in the `#nt-...`
  anchors; F&O regex differences matter for `regex.FO` (see
  docs/LESSONS.md item 6).
