// Package regex implements the two regular-expression flavors XSD 1.1
// requires — one recursive-descent parser, two flavor flags, never one
// flavor for both (docs/LESSONS.md 6):
//
//   - Flavor XSD: the pattern-facet grammar (Datatypes Appendix G).
//     Implicitly anchored; ^ and $ are literal characters; `.` excludes
//     both \n and \r.
//   - Flavor FO: the XQuery/XPath Functions & Operators grammar used by
//     fn:matches / fn:replace / fn:tokenize inside assertions. ^ and $
//     are real anchors; `.` excludes only \n.
//
// The two dot-sets and anchor treatments are the known divergence points;
// any further divergence discovered gets a test naming both flavors.
//
// Translation targets Go regexp/syntax where faithful; constructs Go
// cannot express (e.g. character class subtraction) are rewritten during
// translation. Fuzz targets compare accepted/rejected inputs across
// translation to rule out panics and mistranslation.
//
// Grown during milestone M3 (patterns) and M6 (FO flavor).
package regex
