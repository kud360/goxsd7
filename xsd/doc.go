// Package xsd is the XSD 1.1 schema component model — the in-memory
// representation of a compiled schema (Structures §2.2).
//
// Contract (see docs/ARCHITECTURE.md):
//
//   - PURE LEAF: this package imports nothing from this module except
//     xsderr. Value implementations, regex, parsing, and validation live
//     in packages above it.
//   - Components are immutable after construction and built in phases by
//     package parser (parse → resolve → finalize), so no consumer ever
//     needs a cycle check (docs/STYLE.md D4). Spec-forbidden circularities
//     are rejected during finalization with their named src-/cos- rule.
//   - All child collections are slices in document order; maps are
//     internal lookup indexes only and never determine order (D2).
//   - Nothing derivable is stored: effective facets, transitive
//     membership, etc. are computed on demand (D3; docs/LESSONS.md 5).
//   - Closed sets (variety, derivation method, process contents, …) are
//     typed constants, never strings (STYLE T1).
//
// The component kinds to be modeled (Structures §2.2.1–2.2.3): simple and
// complex type definitions, element and attribute declarations, attribute
// and model groups, particles, wildcards, identity-constraint definitions,
// type alternatives, assertions, notations, and annotations.
//
// Grown during milestone M4 (docs/PLAN.md); QName is provided now as the
// shared naming currency.
package xsd
