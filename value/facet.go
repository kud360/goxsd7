package value

// The facet pipeline (docs/ARCHITECTURE.md) validates a literal against a
// simple type in fixed stages:
//
//	raw literal
//	  → whiteSpace normalization      (lexical)
//	  → LexicalFacet checks           (pattern; every derivation step)
//	  → Mapping.Parse                 (lexical → value)
//	  → ValueFacet checks             (bounds, digits, length, enumeration)
//	  → assertions                    (XPath; fail-open)
//
// Users compose these interfaces to define their own facets and types;
// the builtin facets (Datatypes §4.3) are implemented over the capability
// interfaces in value.go, never over concrete types (STYLE T2). The
// pipeline engine itself lands in milestone M3.

// LexicalFacet constrains the lexical space. Check receives the
// whiteSpace-normalized literal. Errors must be *xsderr.Error with the
// facet's cvc rule (e.g. cvc-pattern-valid).
type LexicalFacet interface {
	CheckLexical(normalized string) error
}

// ValueFacet constrains the value space. Check receives the parsed value.
// Errors must be *xsderr.Error with the facet's cvc rule (e.g.
// cvc-maxInclusive-valid, cvc-length-valid).
type ValueFacet interface {
	CheckValue(v Value) error
}
