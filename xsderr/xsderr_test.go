package xsderr

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorFormatting(t *testing.T) {
	loc := Loc{URI: "schema.xsd", Line: 12, Col: 5}
	err := New("cvc-complex-type.2.1", loc, "element %q must be empty", "foo")
	want := `schema.xsd:12:5: [cvc-complex-type.2.1] element "foo" must be empty`
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestWrapAndExtract(t *testing.T) {
	cause := errors.New("boom")
	inner := Wrap(cause, "src-resolve", Loc{URI: "a.xsd", Line: 3, Col: 1}, "resolving type %q", "t")
	outer := fmt.Errorf("compiling schema: %w", inner)

	if !errors.Is(outer, cause) {
		t.Error("wrapped cause not reachable via errors.Is")
	}
	rule, ok := RuleOf(outer)
	if !ok || rule != "src-resolve" {
		t.Errorf("RuleOf = %q, %v; want src-resolve, true", rule, ok)
	}
	loc, ok := LocOf(outer)
	if !ok || loc.Line != 3 {
		t.Errorf("LocOf = %v, %v; want line 3, true", loc, ok)
	}
}

func TestZeroLoc(t *testing.T) {
	if got := (Loc{}).String(); got != "<unknown>" {
		t.Errorf("zero Loc String() = %q", got)
	}
}
