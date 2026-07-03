package value

import (
	"errors"
	"testing"
)

type fakeMapping string

func (fakeMapping) Parse(string) (Value, error)     { return nil, errors.New("unused") }
func (fakeMapping) Canonical(Value) (string, error) { return "", errors.New("unused") }

type fakeBackend map[string]Mapping

func (b fakeBackend) Mapping(builtin string) (Mapping, bool) {
	m, ok := b[builtin]
	return m, ok
}

func TestOverride(t *testing.T) {
	base := fakeBackend{"decimal": fakeMapping("base-decimal"), "string": fakeMapping("base-string")}
	over := fakeBackend{"decimal": fakeMapping("over-decimal")}
	be := Override(base, over)

	m, ok := be.Mapping("decimal")
	if !ok || m != fakeMapping("over-decimal") {
		t.Errorf("decimal = %v, %v; want override mapping", m, ok)
	}
	m, ok = be.Mapping("string")
	if !ok || m != fakeMapping("base-string") {
		t.Errorf("string = %v, %v; want base mapping", m, ok)
	}
	if _, ok := be.Mapping("gYear"); ok {
		t.Error("gYear should be absent from both")
	}
}
