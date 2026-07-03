// Package validate assesses instance documents against a compiled
// schema, per the cvc-* validation rules of Structures.
//
// Contract:
//
//   - Abstract infoset: instances arrive through small marker interfaces
//     (element/attribute/node views) so formats plug in as adapters —
//     XML now (backed by parser/xmltree), JSON someday. No encoding/xml
//     types in signatures (docs/LESSONS.md 4).
//   - Content-model matching is greedy and deterministic — UPA makes
//     backtracking unnecessary; consume all occurrence matches before
//     exiting a loop, and never let an open-content wildcard absorb what
//     the explicit model can still consume (LESSONS 11).
//   - Known traps threaded from the start: parent element through the
//     assessment chain (cvc-id binds to the parent — LESSONS 14),
//     namespace context through IDC matchers (xpathDefaultNamespace
//     applies to element steps only — LESSONS 12), empty content stricter
//     than element-only (LESSONS 9), EDC via the post-xsi:type governing
//     type (LESSONS 10), union validation against DirectMembers with the
//     validating member's whiteSpace (LESSONS 7).
//   - Every violation is an *xsderr.Error with cvc rule ID and instance
//     (and where relevant schema) Loc. Multiple violations are reported
//     in document order (STYLE D1).
//   - xsi:schemaLocation hints route through the same loader.Resolver the
//     parser uses.
//
// Grown during milestone M5 (docs/PLAN.md).
package validate
