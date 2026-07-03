# goxsd7 Architecture

## Dependency rule

Packages form a strict DAG. `xsd` (the component model) is a **pure leaf**:
it imports nothing from this module. Value implementations, parsing,
validation, and editing live above it. goxsd5 died a little every time this
was violated; see docs/LESSONS.md.

```
                 xsderr          (leaf: errors, rule IDs, locations)
                 xsd             (leaf: component model; imports xsderr only)
                 value           (value-space contracts, facet pipeline; imports xsd, xsderr)
   builtin/strict  builtin/native  <user backends>   (implement value contracts)
                 regex           (XSD + F&O flavors)
                 parser/xmltree  (position-tracking XML; independent)
                 loader          (schema resolution interfaces)
                 parser          (schema docs -> xsd components)
                 xpath           (subset engine; imports value)
                 validate        (instance validation)
                 codegen  codec  (generation; dataset ser/de)
                 conformance     (harness + ratchet; test-only)
```

## Lexical space vs value space

The load-bearing separation of the whole design (Datatypes §2.1–2.3):

- **Lexical space**: strings. Whitespace normalization and `pattern` facets
  operate here, *before* any parsing.
- **Value space**: typed Go values. Ordering, equality, identity, and all
  value-based facets (`minInclusive`, `totalDigits`, `length` on lists,
  `enumeration`, `assertion`) operate here, *after* the lexical mapping.
- The bridge is the pair of **mappings** per type: lexical → value and
  value → canonical lexical. These are defined normatively as function
  definitions ("hfn") in Datatypes Appendix E, and our builtins are
  bootstrapped from those definitions (extracted from
  `docs/specs/md/xmlschema11-2.md` by `tools/hfnextract`), not
  hand-transcribed.

### The facet pipeline

Validation of a literal against a simple type is a fixed pipeline; each
stage is a value users can compose for their own types:

```
raw literal
  → whiteSpace normalization        (lexical; from the type's ws facet)
  → pattern facets                  (lexical; every step of the derivation chain)
  → lexical mapping                 (string → value.Value, via the backend)
  → value facets                    (bounds, digits, length, enumeration)
  → assertions                      (XPath, fail-open; per-item for lists,
                                     per-member for unions, at every level)
```

List and union varieties recurse: lists apply the pipeline per item against
the item type before list-level facets; unions try DirectMembers in order
(not flattened members — intervening restrictions carry facets, and pattern
normalization uses the *validating member's* whiteSpace).

## Values, capabilities, and user backends

`value.Value` is `any`. Capabilities are small interfaces the facet pipeline
and comparison logic detect:

- `value.Eq` / `value.Ordered` — equality and the partial order
- `value.Lengthed` — length-facet units (chars, octets, list items)
- `value.DigitCounted` — totalDigits / fractionDigits
- `value.TimezoneAware` — date/time timezone handling
- `value.Canonical` — canonical lexical form

No sealed interfaces, no exhaustive type switches outside the defining
package: that is what lets **users bring their own backend**.

### Backends are pluggable, per type, with fallback

A `value.Backend` maps builtin type names to `value.Mapping` (the
lexical→value and value→canonical pair plus capability wiring). Ships with:

- `builtin/strict` — spec-exact: arbitrary-precision decimal/integer,
  `precisionDecimal`, the 7-property date/time model, XSD-exact
  float/double behavior.
- `builtin/native` — Go-friendly: `int64`, `float64`, `string`,
  `time.Time`; documented, deliberate deviations from the spec value
  spaces (range limits, timezone folding).

Backends compose: `value.Override(base, partial)` produces a backend that
uses `partial`'s mappings where defined and falls back to `base`. So a user
can back only `xs:decimal` with their money type and keep strict everything
else — or implement a full backend of their own. User-defined *types*
(restrictions/lists/unions over builtins) get the pipeline for free; only
primitive mappings are backend territory.

## Component model (`xsd`)

- Components are constructed in **phases** so no traversal ever needs a
  cycle check (STYLE D4): (1) parse schema documents into raw form,
  (2) resolve QName references through a symbol table,
  (3) finalize in dependency order — a component's base/item/member types
  are complete before it is. Spec-forbidden circularities (`st-props-correct`
  circular unions, circular substitution groups, …) are rejected at phase 3
  with their named rule.
- All child collections are slices in document order. Maps exist only as
  internal indexes and never determine any order.
- Nothing derivable is stored (no effective-facet caches: compute
  `Merge(base.EffectiveFacets(), declared)` on demand).
- The model is **read-only** after construction; mutation/editing APIs are
  out of scope for now.

## Parsing & loading

- `parser/xmltree`: streaming, bounded-memory XML reader that records
  line/column for every node; the origin of every `xsderr.Loc`. No
  `io.ReadAll`.
- `loader`: the IO seam. `Resolver` answers "give me the schema document
  for (namespace, location hint)"; helpers provided for files, HTTP, and
  in-memory maps, plus a chaining/catalog resolver. `xsi:schemaLocation`
  instance hints route through the same interface so multi-schema loading
  stays in one place (a goxsd5 deferred problem — design for it now).

## XPath (`xpath`)

Grown from the XSD-required subset outward: CTA's restricted subset first,
then assertion essentials (axes, predicates, quantified expressions, F&O
functions). **Fail-open**: an unsupported construct can never cause a false
rejection; every fallback site is a greppable `// GAP(xpath): …`. F&O regex
functions use the F&O flavor (`regex.FO`), never the pattern-facet flavor —
the flavors differ on anchors and dot semantics (see LESSONS).

## Validation (`validate`)

- Abstract infoset via marker interfaces so XML today / JSON later plug in
  as adapters.
- Content-model matching is greedy and deterministic (UPA makes
  backtracking unnecessary); explicit content beats open-content wildcards.
- Every violation is an `xsderr.Error` with rule ID + instance and/or
  schema location.

## Codegen & codec

- `codegen` emits Go types from a compiled schema, deterministically
  (D1/D2). Users choose the strict or native backend for generated fields.
- **Choices are sealed interfaces.** An `xs:choice` becomes an interface
  with an unexported marker method; each branch is a concrete type
  implementing it, and consumers use type switches. This is the closed-sum
  exception to STYLE T2 (spelled out there): exactly one branch can exist,
  so "N pointer fields, exactly one non-nil" — an illegal-states factory —
  never appears in generated code.
- **Anonymous types get ancestor-context names.** A single namer component
  owns all XSD-name → Go-identifier decisions. Anonymous types are named
  by walking up their schema ancestors to the nearest named declaration
  (element `shipTo` under element `purchaseOrder` → `PurchaseOrderShipTo`),
  extending the path only as far as uniqueness requires; residual
  collisions (case folding, Go keywords, XML-legal-but-Go-illegal names)
  are disambiguated deterministically by document order (D1/D2). Every
  generated type's header comment records its schema Loc + original QName.
- `codec` is the dataset serializer/deserializer: schema-directed decode of
  instance documents into generated (or reflective) Go values and canonical
  encode back out.

### Two decode paths, one semantics

`codec` is built for **minimal allocation**:

- **Runtime path** (always available): the facet pipeline +
  `value.Mapping`, driven by the compiled schema. General, reflective,
  allocation-tolerant.
- **Generated fast path**: builtin backends export their own **code
  emitters** (`value.Emitter`, implemented by `builtin/strict` and
  `builtin/native`; user backends may implement it too). At codegen time
  the emitter contributes specialized decode/encode code for its types —
  parsing directly from the reader's byte window into the target field,
  no intermediate string, no boxed `value.Value`, facet checks inlined.
  A backend without an emitter simply falls back to the runtime path for
  its types.
- Runtime hot-path APIs follow the appender convention
  (`AppendCanonical(dst []byte, v) []byte`, `ParseBytes(b []byte)`) so
  even the non-generated path can be allocation-frugal.

The two paths implement the *same* pipeline stages with the *same* spec
rule IDs, which makes them **differentially testable**: for every type,
property tests feed identical input to both paths and require identical
values and identical error rule IDs. A fast path that disagrees with the
runtime path is wrong by definition — this is the primary defense against
"optimized but subtly different" parsing.

### Debuggability of parsing

When a value fails to parse, the error must localize the failure without
a debugger (extending E1–E3):

- every decode error carries the **pipeline stage** that rejected
  (whitespace / pattern / lexical-map / facet / assertion), the type
  QName, the offending input fragment, and the instance Loc + byte offset;
- `GOXSD_DEBUG=codec` traces stage transitions per value through the
  injected slog logger (rule ID, type, input) so an agent can watch one
  value flow through the pipeline;
- generated code preserves this: emitted fast paths report the same
  stage/rule metadata as the runtime path, and generated files map cleanly
  back to the emitting backend and schema construct (a header comment per
  emitted decode function naming type QName + schema Loc).

## Conformance & ratchet

- W3C suite at `testdata/xsdtests` (submodule, pinned).
- Expectations committed at `conformance/testdata/expectations/*.txt`, one
  line per test case; diffs make regressions obvious and `git blame`
  bisectable.
- `make conformance` compares; `make ratchet` re-baselines **upward only**.
  A regression fails loudly and must never be committed.
- XPath gets its own ratchet lane once the engine grows past the subset
  (W3C XPath tests, `conformance/xpath`).

## Logging

`log/slog` injected at construction, namespaced groups, silent by default.
The debug level is designed for agents: messages carry rule ID, component
QName, and location so a conformance failure can be localized from logs
alone (`GOXSD_DEBUG=parser,validate` in tests).
