package conformance

import (
	"os"
	"path/filepath"
	"testing"
)

// suiteDir is the W3C test suite submodule; suite.xml is its index of
// test-set references.
const suiteDir = "../testdata/xsdtests"

// TestConformance is the ratchet entry point (`make conformance` /
// `make ratchet`). Lanes come online per docs/PLAN.md: datatypes (M3),
// schema validity (M4), instance validity (M5), xpath (M6+). Until the
// first processor slice exists there is nothing to run — but the suite
// and expectations layout are verified so the harness fails loudly if
// the environment is broken.
func TestConformance(t *testing.T) {
	if _, err := os.Stat(filepath.Join(suiteDir, "suite.xml")); err != nil {
		t.Fatalf("W3C suite not available (run `git submodule update --init`): %v", err)
	}

	lanes, err := filepath.Glob("testdata/expectations/*.txt")
	if err != nil {
		t.Fatalf("globbing expectations: %v", err)
	}
	for _, lane := range lanes {
		if _, err := LoadExpectations(lane); err != nil {
			t.Errorf("lane %s is corrupt: %v", lane, err)
		}
	}

	// GAP(conformance): no processor lanes implemented yet (PLAN M3).
	// The runner will: parse suite.xml → test sets → for each case run
	// the processor, derive Pass/Fail vs the suite's declared expected
	// outcome, Compare against the lane file, fail on Delta.Regressed,
	// and under GOXSD_RATCHET=1 call Ratchet and rewrite the lane.
	t.Skip("no conformance lanes implemented yet (docs/PLAN.md M3)")
}
