// Package native is the Go-friendly value backend: builtin datatypes
// mapped to natural Go representations, trading spec-exactness for
// ergonomics.
//
//   - integer family → int64 (out-of-range literals are errors here,
//     valid in builtin/strict)
//   - decimal → deliberate approximation (documented), float/double → float64
//   - date/time family → time.Time (timezone-less values get a documented
//     folding; year/leap-second edge cases follow Go, not XSD)
//   - string family → string
//
// Every deviation from the spec value spaces is enumerated in this
// package's deviation table (deviations.md, milestone M9) — a deviation
// that isn't documented is a bug.
//
// Backend() returns the value.Backend. Mix-and-match with strict or user
// backends via value.Override: users may take any subset of these types
// and supply their own for the rest.
//
// Like builtin/strict, this package exports ParseBytes/AppendCanonical
// hot-path entry points and its own codegen emitter (value.Emitter, M7);
// native is the backend where scalar decode on the generated fast path
// targets zero allocations (int64/float64/bool straight from the byte
// window).
//
// Grown during milestone M9 (docs/PLAN.md).
package native
