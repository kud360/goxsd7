package conformance

import (
	"strings"
	"testing"
)

func TestParseAndWriteRoundTrip(t *testing.T) {
	in := "# comment\n\nb-case fail\na-case pass\n"
	res, err := ParseExpectations(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseExpectations: %v", err)
	}
	var out strings.Builder
	if err := WriteExpectations(&out, res); err != nil {
		t.Fatalf("WriteExpectations: %v", err)
	}
	want := "a-case pass\nb-case fail\n" // sorted, comments dropped
	if out.String() != want {
		t.Errorf("round trip = %q, want %q", out.String(), want)
	}
}

func TestParseRejectsBadInput(t *testing.T) {
	for name, in := range map[string]string{
		"missing status": "case1\n",
		"bad status":     "case1 maybe\n",
		"duplicate":      "case1 pass\ncase1 fail\n",
	} {
		if _, err := ParseExpectations(strings.NewReader(in)); err == nil {
			t.Errorf("%s: expected parse error, got none", name)
		}
	}
}

func TestCompare(t *testing.T) {
	expected := Results{"a": Pass, "b": Fail, "c": Pass, "d": Fail}
	actual := Results{"a": Pass, "b": Pass, "c": Fail, "e": Fail}
	d := Compare(expected, actual)

	check := func(name string, got []string, want ...string) {
		t.Helper()
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Errorf("%s = %v, want %v", name, got, want)
		}
	}
	check("Improved", d.Improved, "b")
	check("Regressed", d.Regressed, "c")
	check("New", d.New, "e")
	check("Vanished", d.Vanished, "d")
}

func TestRatchetRefusesRegression(t *testing.T) {
	expected := Results{"a": Pass}
	actual := Results{"a": Fail}
	if _, err := Ratchet(expected, actual); err == nil {
		t.Fatal("Ratchet accepted a regression")
	}
}

func TestRatchetMergesUpward(t *testing.T) {
	expected := Results{"a": Pass, "b": Fail}
	actual := Results{"a": Pass, "b": Pass, "c": Fail}
	merged, err := Ratchet(expected, actual)
	if err != nil {
		t.Fatalf("Ratchet: %v", err)
	}
	if merged["b"] != Pass {
		t.Error("improved case b not harvested")
	}
	if merged["c"] != Fail {
		t.Error("new case c not recorded")
	}
	passed, total := Counts(merged)
	if passed != 2 || total != 3 {
		t.Errorf("Counts = %d/%d, want 2/3", passed, total)
	}
}
