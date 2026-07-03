// Package value defines the value-space contracts of goxsd7: what a typed
// value is, which capabilities it may offer, and how lexical↔value
// mappings are supplied by pluggable backends.
//
// Lexical space vs value space (Datatypes §2.1–2.3) is the load-bearing
// separation: strings and pattern facets live in the lexical space; typed
// values, ordering, and value-based facets live here. The bridge is a
// Mapping per primitive type, defined normatively by the function
// definitions in Datatypes Appendix E (anchor ap-funcDefs, functions
// f-*Lexmap / f-*Canmap) — builtins are bootstrapped from those, not
// hand-transcribed (tools/hfnextract, milestone M1).
//
// Capabilities are interfaces, never type switches (docs/STYLE.md T2):
// that is what lets users bring their own value backend, wholesale or for
// any subset of builtin types (Override).
package value

// Value is any value-space datum. Deliberately open (docs/LESSONS.md 2):
// user backends provide their own concrete types and participate by
// implementing the capability interfaces below.
type Value = any

// Ordering is the outcome of comparing two values in a value space.
// XSD value spaces are partially ordered: values from different primitive
// spaces — and some within one space, e.g. dateTimes with and without
// timezone — are Incomparable (Datatypes §4.2.1).
type Ordering int

const (
	Less         Ordering = iota - 1
	Equal                 // also: identity for facet purposes where specified
	Greater               //
	Incomparable          //
)

func (o Ordering) String() string {
	switch o {
	case Less:
		return "less"
	case Equal:
		return "equal"
	case Greater:
		return "greater"
	}
	return "incomparable"
}

// Eq is the equality capability (Datatypes §2.2.1). Required of every
// value: enumeration, fixed values, identity constraints all need it.
type Eq interface {
	Eq(other Value) bool
}

// Ordered adds the (partial) order capability; needed by the min/max
// bounding facets.
type Ordered interface {
	Eq
	Cmp(other Value) Ordering
}

// Lengthed reports length in the unit the type measures (characters for
// strings, octets for binary, items for lists); needed by length,
// minLength, maxLength.
type Lengthed interface {
	Len() int
}

// DigitCounted exposes digit counts for totalDigits / fractionDigits.
type DigitCounted interface {
	TotalDigits() int
	FractionDigits() int
}

// TimezoneAware distinguishes date/time values with and without timezone
// offsets (the explicitTimezone facet, and Incomparable ordering cases).
type TimezoneAware interface {
	HasTimezone() bool
}

// Canonical renders the canonical lexical representation (the value →
// literal mapping of Appendix E). Values lacking it cannot be
// canonically serialized by codec.
type Canonical interface {
	Canonical() string
}
