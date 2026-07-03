# Ratchet expectation lanes

One file per lane, committed, sorted, machine-written (`make ratchet`) —
never hand-edited. Format: `<case-id> <pass|fail>` per line, `#` comments.

Planned lanes (created by the first ratchet run of each):

- `datatypes.txt` — simple-type conformance (from PLAN M3)
- `schema.txt` — schema validity across the full suite (M4)
- `instance.txt` — instance validity (M5)
- `xpath.txt` — XPath lane (M6+)

Rules: upward only; every flip explainable by the diff that caused it;
regressions block commit (see docs/WORKFLOW.md and conformance/ratchet.go).
