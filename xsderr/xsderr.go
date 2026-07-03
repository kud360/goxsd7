// Package xsderr defines the error currency of goxsd7.
//
// Every schema- or instance-validity violation is an *Error carrying the
// spec validation rule it violates (Rule), the source location it arose at
// (Loc), and human context. See docs/STYLE.md E1–E3: errors are always
// decorated, always rule-mapped, always located.
//
// This package is a pure leaf: it imports only the standard library.
package xsderr

import (
	"errors"
	"fmt"
)

// Rule identifies a numbered constraint in the XSD 1.1 specifications,
// e.g. "cvc-complex-type.2.1", "cos-st-restricts.1.1", "src-resolve".
// The rule IDs are anchors in docs/specs/md and must be quoted verbatim.
type Rule string

func (r Rule) String() string { return string(r) }

// Loc is a position in a source document (schema or instance).
// The zero Loc means "location unknown" and should be rare: parsers must
// thread positions (docs/STYLE.md E3).
type Loc struct {
	URI  string // document URI or path
	Line int    // 1-based; 0 = unknown
	Col  int    // 1-based; 0 = unknown
}

func (l Loc) IsZero() bool { return l == Loc{} }

func (l Loc) String() string {
	if l.IsZero() {
		return "<unknown>"
	}
	return fmt.Sprintf("%s:%d:%d", l.URI, l.Line, l.Col)
}

// Error is a spec-rule violation. Treat as immutable after construction.
type Error struct {
	Rule Rule
	Loc  Loc
	Msg  string
	Err  error // wrapped cause; may be nil
}

// New creates a rule violation at loc.
func New(rule Rule, loc Loc, format string, args ...any) *Error {
	return &Error{Rule: rule, Loc: loc, Msg: fmt.Sprintf(format, args...)}
}

// Wrap decorates an underlying cause with rule and location.
func Wrap(err error, rule Rule, loc Loc, format string, args ...any) *Error {
	return &Error{Rule: rule, Loc: loc, Msg: fmt.Sprintf(format, args...), Err: err}
}

func (e *Error) Error() string {
	s := fmt.Sprintf("%s: [%s] %s", e.Loc, e.Rule, e.Msg)
	if e.Err != nil {
		s += ": " + e.Err.Error()
	}
	return s
}

func (e *Error) Unwrap() error { return e.Err }

// RuleOf extracts the spec rule from err, if it (or anything it wraps)
// is an *Error.
func RuleOf(err error) (Rule, bool) {
	var e *Error
	if !errors.As(err, &e) {
		return "", false
	}
	return e.Rule, true
}

// LocOf extracts the nearest source location from err.
func LocOf(err error) (Loc, bool) {
	var e *Error
	if !errors.As(err, &e) {
		return Loc{}, false
	}
	return e.Loc, true
}
