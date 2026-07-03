package value

// Mapping binds one builtin datatype to its two normative mappings
// (Datatypes Appendix E): lexical → value and value → canonical literal.
// Parse receives the literal already whiteSpace-normalized and
// pattern-checked (the lexical stages of the facet pipeline run first;
// see docs/ARCHITECTURE.md).
type Mapping interface {
	Parse(lexical string) (Value, error)
	Canonical(v Value) (string, error)
}

// Backend supplies Mappings for builtin datatypes, keyed by local name in
// the XSD namespace (e.g. "decimal", "dateTime"). goxsd7 ships
// builtin/strict (spec-exact) and builtin/native (Go-friendly); users may
// implement Backend themselves — returning ok=false for types they do not
// cover — and compose with Override.
//
// User-defined *types* (restrictions, lists, unions over builtins) are
// not backend territory: they get the facet pipeline for free; only
// primitive mappings vary by backend.
type Backend interface {
	Mapping(builtin string) (Mapping, bool)
}

// Backends may additionally implement the emitter seam (value.Emitter,
// API fixed in milestone M7): a codegen-time contract through which a
// backend contributes specialized zero-allocation decode/encode code for
// its types. Emitting is optional — types whose backend has no emitter
// use the codec runtime path; both paths must agree (differential tests,
// see codec doc).

// Override composes two backends: mappings defined by over win, base
// fills the rest. Typical use — back only xs:decimal with a custom money
// type and keep the strict backend for everything else:
//
//	be := value.Override(strict.Backend(), myMoneyBackend)
func Override(base, over Backend) Backend {
	return overrideBackend{base: base, over: over}
}

type overrideBackend struct {
	base Backend
	over Backend
}

func (o overrideBackend) Mapping(builtin string) (Mapping, bool) {
	if m, ok := o.over.Mapping(builtin); ok {
		return m, true
	}
	return o.base.Mapping(builtin)
}
