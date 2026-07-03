// Package xpath is the XPath 2.0 engine, grown outward from the subset
// XSD 1.1 requires:
//
//  1. CTA first — the restricted subset for conditional type alternatives
//     (Structures, test attribute of xs:alternative).
//  2. Assertion essentials — axes, predicates, quantified expressions,
//     typed comparisons, the F&O function core. fn:matches/replace/
//     tokenize bind to regex flavor FO, never the pattern flavor
//     (docs/LESSONS.md 6).
//  3. Long-term: the full grammar (docs/specs/md/xpath20.md), with an
//     XPath conformance lane in the ratchet.
//
// FAIL-OPEN CONTRACT (docs/LESSONS.md 16): an unsupported construct must
// never cause a false rejection — assertions evaluate as satisfied, type
// alternatives as unmatched. Every fail-open site carries a
// `// GAP(xpath): <construct>` comment; the cartographer turns those into
// kind/gap issues, and the ratchet closes them over time. A fail-open
// without a GAP marker is a style violation (STYLE P3).
//
// Variables are typed, not stringly: $value binds a typed atom carrying
// {Lexical, Kind} so comparison semantics are correct (LESSONS 15).
// Assertions apply at every level of the variety chain — per-item for
// lists, per-member for unions (LESSONS 8).
//
// Grown during milestone M6 (docs/PLAN.md).
package xpath
