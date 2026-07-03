// Package codegen emits Go source from a compiled schema: struct types
// for complex types, named types for simple types, with users choosing
// the builtin/strict or builtin/native backend (or their own) for value
// fields.
//
// Contract:
//
//   - Deterministic output (STYLE D1/D2): identical schema in, byte-
//     identical Go out. Declarations in schema document order; no map
//     iteration anywhere near the emitter.
//   - Generated code obeys docs/STYLE.md like handwritten code, and
//     compiles under `go vet` cleanly — golden-file tests enforce both.
//   - CHOICES ARE SEALED INTERFACES (the closed-sum exception in STYLE
//     T2): an xs:choice becomes an interface with an unexported marker
//     method, one concrete branch type per alternative, consumers type-
//     switch. Never N pointer fields with "exactly one non-nil".
//     Required/optional and nillable likewise map to shapes that cannot
//     express the invalid combination (T1) — the warden reviews the
//     generated shapes, not just the generator.
//   - FAST-PATH EMISSION: value backends that implement the emitter seam
//     (value.Emitter, API fixed in milestone M7) contribute specialized
//     decode/encode code for their builtin types — parsing straight from
//     the reader's byte window into the target field, facet checks
//     inlined, no boxed value.Value. Types whose backend has no emitter
//     fall back to the codec runtime path; both paths carry identical
//     spec rule IDs and are differentially tested (see codec doc).
//   - Generated decode/encode functions are traceable back to their
//     origin: a header comment per function naming the type QName, the
//     schema Loc, and the emitting backend.
//   - NAMING: one namer component owns every XSD-name → Go-identifier
//     decision. Anonymous types are named from their schema ancestor
//     context (enclosing element/attribute declarations up to the nearest
//     named ancestor: element "shipTo" inside element "purchaseOrder" →
//     PurchaseOrderShipTo), lengthening the ancestor path only as far as
//     needed for uniqueness. Name wrangling (case folding, separators,
//     Go keyword and predeclared-identifier collisions, XSD names legal
//     in XML but not in Go) resolves conflicts deterministically —
//     collisions get a stable, document-order-based disambiguation, never
//     a random or map-order one (D1/D2). The chosen name is recorded in
//     the generated type's header comment next to its schema Loc so
//     users can trace a Go type back to the schema construct.
//   - User-defined derived types keep working: generated types for
//     restrictions expose their facet pipeline so runtime validation is
//     available, not compiled away.
//
// Grown during milestone M7 (docs/PLAN.md).
package codegen
