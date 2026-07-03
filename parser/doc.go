// Package parser compiles schema documents into xsd components.
//
// Construction is phased (docs/STYLE.md D4) so the finished component
// graph never needs cycle checks:
//
//	phase 1  parse    — each schema document → raw form (via
//	                    parser/xmltree; every raw node keeps its Loc)
//	phase 2  resolve  — schema composition (include/import/redefine/
//	                    override, chameleon namespace coercion) and QName
//	                    reference resolution through a symbol table
//	phase 3  finalize — components completed in dependency order; a
//	                    component's base/item/member types are finished
//	                    before it is. Spec-forbidden circularities
//	                    (st-props-correct circular unions, circular
//	                    groups, …) are rejected HERE, once, with their
//	                    named src-/cos- rule and location.
//
// xs:override requires explicit target tracking — schema-level defaults
// and groups declared inside the override attach to the overridden
// document (docs/LESSONS.md 13).
//
// All schema loading goes through a loader.Resolver. Every error is an
// *xsderr.Error with rule ID and schema Loc.
//
// Grown during milestone M4 (docs/PLAN.md).
package parser
