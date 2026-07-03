// Package strict is the spec-exact value backend: every builtin datatype
// represented with full fidelity to its value space in Datatypes §3 and
// the precisionDecimal spec.
//
//   - decimal/integer families: arbitrary precision (math/big based),
//     digit counts preserved for totalDigits/fractionDigits.
//   - precisionDecimal: IEEE 754 decimal semantics per
//     docs/specs/md/xsd-precisionDecimal.md (NaN, ±INF, signed zero,
//     precision-carrying values).
//   - date/time family: the seven-property model (Datatypes §D.2.1,
//     anchor theSevenPropertyModel) — NOT time.Time; year 0 handling,
//     timezone-less values, and ±14:00 offsets are spec-exact.
//   - float/double: XSD-exact behavior including NaN equality-to-itself
//     for value identity purposes.
//
// The mappings are bootstrapped from the Appendix E function definitions
// extracted by tools/hfnextract (milestone M1); parsing/canonical code is
// generated or table-driven from those definitions, not hand-transcribed
// (milestone M3).
//
// Backend() returns the value.Backend covering all builtin types. Users
// needing Go-friendly representations use builtin/native or their own
// backend, composed via value.Override.
//
// Hot-path surface: parsing exposes ParseBytes-style entry points and
// canonicalization the appender convention (AppendCanonical(dst []byte,
// v) []byte). The package exports its own codegen emitter (value.Emitter,
// M7) so generated code can decode strict values straight from reader
// byte windows with facet checks inlined — differential tests pin the
// emitter's semantics to the runtime mappings.
package strict
