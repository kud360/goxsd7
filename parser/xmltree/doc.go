// Package xmltree reads XML documents into a position-annotated tree.
// It is the origin of every xsderr.Loc: each element, attribute, and
// character-data node records its line and column.
//
// Contract (docs/LESSONS.md 21):
//
//   - Streaming input, bounded buffering — never io.ReadAll. Peak memory
//     is O(returned tree), not O(buffers).
//   - Line/column via an offset index, not stored line starts.
//   - Namespace resolution performed here; consumers see expanded names.
//   - Encodings and the XML prolog handled here (only the DTD internal
//     subset may be buffered).
//   - Independent leaf package: usable by parser (schema documents) and
//     validate (instance documents) alike.
//
// Grown during milestone M2 (docs/PLAN.md), with a fuzz target from the
// start (fuzzing rules out panics on malformed input — LESSONS 20).
package xmltree
