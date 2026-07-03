// Package conformance runs the W3C XML Schema test suite
// (testdata/xsdtests, git submodule) against goxsd7 and enforces the
// ratchet: committed expectations only ever improve.
//
// Expectations live in testdata/expectations/*.txt, one lane per file
// (schema validity, instance validity, later xpath). Format: one case per
// line, `<case-id> <status>`, sorted by case ID; '#' starts a comment.
// Status is pass or fail — "fail" cases are known gaps, recorded so any
// regression of a "pass" case is loud and any improvement is harvestable
// via `make ratchet` (GOXSD_RATCHET=1).
//
// Rules (docs/WORKFLOW.md): expectations move upward only, are never
// hand-edited, and every flip must be explainable by the diff that
// caused it.
package conformance

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
)

// Status is a case outcome relative to the suite's declared expectation:
// Pass means goxsd7 agrees with the suite, Fail means it does not (yet).
type Status string

const (
	Pass Status = "pass"
	Fail Status = "fail"
)

func parseStatus(s string) (Status, error) {
	switch Status(s) {
	case Pass, Fail:
		return Status(s), nil
	}
	return "", fmt.Errorf("unknown status %q (want pass|fail)", s)
}

// Results maps case ID → status for one lane. Order never leaves this
// map: all output goes through sorted case IDs (STYLE D2).
type Results map[string]Status

// ParseExpectations reads a lane file.
func ParseExpectations(r io.Reader) (Results, error) {
	res := Results{}
	sc := bufio.NewScanner(r)
	line := 0
	for sc.Scan() {
		line++
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		id, statusStr, ok := strings.Cut(text, " ")
		if !ok {
			return nil, fmt.Errorf("line %d: want `<case-id> <status>`, got %q", line, text)
		}
		status, err := parseStatus(strings.TrimSpace(statusStr))
		if err != nil {
			return nil, fmt.Errorf("line %d (case %s): %w", line, id, err)
		}
		if _, dup := res[id]; dup {
			return nil, fmt.Errorf("line %d: duplicate case %q", line, id)
		}
		res[id] = status
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("reading expectations: %w", err)
	}
	return res, nil
}

// LoadExpectations reads a lane file from disk. A missing file is an
// empty lane (the lane hasn't been ratcheted yet), not an error.
func LoadExpectations(path string) (Results, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return Results{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("opening expectations %s: %w", path, err)
	}
	defer f.Close()
	res, err := ParseExpectations(f)
	if err != nil {
		return nil, fmt.Errorf("parsing expectations %s: %w", path, err)
	}
	return res, nil
}

// WriteExpectations writes a lane deterministically: sorted case IDs,
// one per line.
func WriteExpectations(w io.Writer, res Results) error {
	ids := make([]string, 0, len(res))
	for id := range res {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		if _, err := fmt.Fprintf(w, "%s %s\n", id, res[id]); err != nil {
			return fmt.Errorf("writing expectation for %s: %w", id, err)
		}
	}
	return nil
}

// Delta is the comparison of an actual run against expectations.
// All slices are sorted by case ID.
type Delta struct {
	Improved  []string // expected fail, now pass — harvest via Ratchet
	Regressed []string // expected pass, now fail — NEVER acceptable
	New       []string // cases with no expectation yet
	Vanished  []string // expected cases the run no longer produced
}

// Compare diffs actual results against expectations.
func Compare(expected, actual Results) Delta {
	var d Delta
	for id, got := range actual {
		want, known := expected[id]
		if !known {
			d.New = append(d.New, id)
			continue
		}
		if want == Fail && got == Pass {
			d.Improved = append(d.Improved, id)
		}
		if want == Pass && got == Fail {
			d.Regressed = append(d.Regressed, id)
		}
	}
	for id := range expected {
		if _, ok := actual[id]; !ok {
			d.Vanished = append(d.Vanished, id)
		}
	}
	slices.Sort(d.Improved)
	slices.Sort(d.Regressed)
	slices.Sort(d.New)
	slices.Sort(d.Vanished)
	return d
}

// Ratchet merges an actual run into expectations, upward only. It returns
// an error — and merges nothing — if any case regressed or vanished:
// a ratchet must never record a downgrade (docs/WORKFLOW.md).
func Ratchet(expected, actual Results) (Results, error) {
	d := Compare(expected, actual)
	if len(d.Regressed) > 0 {
		return nil, fmt.Errorf("refusing to ratchet: %d regression(s): %s",
			len(d.Regressed), strings.Join(d.Regressed, ", "))
	}
	if len(d.Vanished) > 0 {
		return nil, fmt.Errorf("refusing to ratchet: %d expected case(s) missing from the run: %s",
			len(d.Vanished), strings.Join(d.Vanished, ", "))
	}
	merged := make(Results, len(expected)+len(d.New))
	maps.Copy(merged, expected)
	for _, id := range d.Improved {
		merged[id] = Pass
	}
	for _, id := range d.New {
		merged[id] = actual[id]
	}
	return merged, nil
}

// Counts summarizes a lane as (passed, total).
func Counts(res Results) (passed, total int) {
	for _, s := range res {
		total++
		if s == Pass {
			passed++
		}
	}
	return passed, total
}
