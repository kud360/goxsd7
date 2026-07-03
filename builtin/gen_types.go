package builtin

// Variety represents whether a type is atomic or a list.
type Variety string

const (
	VarietyAtomic Variety = "atomic"
	VarietyList   Variety = "list"
)

// Facet represents a constraint applied to a datatype.
type Facet struct {
	Name  string
	Value string
}

// BuiltinType defines the properties of an XSD 1.1 built-in datatype.
type BuiltinType struct {
	Name     string
	BaseType string
	Variety  Variety
	Facets   []Facet
}

// BuiltinTypes contains all built-in XSD 1.1 datatypes in hierarchical order.
var BuiltinTypes = []BuiltinType{
	{
		Name:     "anySimpleType",
		BaseType: "",
		Variety:  VarietyAtomic,
	},
	{
		Name:     "anyAtomicType",
		BaseType: "anySimpleType",
		Variety:  VarietyAtomic,
	},

	// Primitive Datatypes (Base: anyAtomicType)
	{
		Name:     "string",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "preserve"},
		},
	},
	{
		Name:     "boolean",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "decimal",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "float",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "double",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "dateTime",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "duration",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "time",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "date",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "gYear",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "gMonth",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "gDay",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "gYearMonth",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "gMonthDay",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "hexBinary",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "base64Binary",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "anyURI",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "QName",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "NOTATION",
		BaseType: "anyAtomicType",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},

	// Ordinary Built-in Datatypes
	{
		Name:     "normalizedString",
		BaseType: "string",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "replace"},
		},
	},
	{
		Name:     "token",
		BaseType: "normalizedString",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "language",
		BaseType: "token",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"pattern", `[a-zA-Z]{1,8}(-[a-zA-Z0-9]{1,8})*`},
		},
	},
	{
		Name:     "NMTOKEN",
		BaseType: "token",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			// Simplified XSD regex for NMTOKEN (excluding colon)
			{"pattern", `[a-zA-Z0-9._\x80-\uFFFF-]+`},
		},
	},
	{
		Name:     "NMTOKENS",
		BaseType: "anySimpleType",
		Variety:  VarietyList,
		Facets: []Facet{
			{"item", "NMTOKEN"},
			{"minLength", "1"},
		},
	},
	{
		Name:     "Name",
		BaseType: "token",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			// Simplified XSD regex for Name (including colon)
			{"pattern", `[a-zA-Z_:\x80-\uFFFF][a-zA-Z0-9._:\x80-\uFFFF-]*`},
		},
	},
	{
		Name:     "NCName",
		BaseType: "Name",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			// Simplified XSD regex for NCName (excluding colon)
			{"pattern", `[a-zA-Z_\x80-\uFFFF][a-zA-Z0-9._\x80-\uFFFF-]*`},
		},
	},
	{
		Name:     "ID",
		BaseType: "NCName",
		Variety:  VarietyAtomic,
		Facets:   []Facet{},
	},
	{
		Name:     "IDREF",
		BaseType: "NCName",
		Variety:  VarietyAtomic,
		Facets:   []Facet{},
	},
	{
		Name:     "IDREFS",
		BaseType: "anySimpleType",
		Variety:  VarietyList,
		Facets: []Facet{
			{"item", "IDREF"},
			{"minLength", "1"},
		},
	},
	{
		Name:     "ENTITY",
		BaseType: "NCName",
		Variety:  VarietyAtomic,
		Facets:   []Facet{},
	},
	{
		Name:     "ENTITIES",
		BaseType: "anySimpleType",
		Variety:  VarietyList,
		Facets: []Facet{
			{"item", "ENTITY"},
			{"minLength", "1"},
		},
	},
	{
		Name:     "integer",
		BaseType: "decimal",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"fractionDigits", "0"},
			{"pattern", `[\-+]?[0-9]+`},
		},
	},
	{
		Name:     "nonPositiveInt",
		BaseType: "integer",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"maxInclusive", "0"},
		},
	},
	{
		Name:     "negativeInteger",
		BaseType: "nonPositiveInt",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"maxInclusive", "-1"},
		},
	},
	{
		Name:     "long",
		BaseType: "integer",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"minInclusive", "-9223372036854775808"},
			{"maxInclusive", "9223372036854775807"},
		},
	},
	{
		Name:     "int",
		BaseType: "long",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"minInclusive", "-2147483648"},
			{"maxInclusive", "2147483647"},
		},
	},
	{
		Name:     "short",
		BaseType: "int",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"minInclusive", "-32768"},
			{"maxInclusive", "32767"},
		},
	},
	{
		Name:     "byte",
		BaseType: "short",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"minInclusive", "-128"},
			{"maxInclusive", "127"},
		},
	},
	{
		Name:     "nonNegativeInt",
		BaseType: "integer",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"minInclusive", "0"},
		},
	},
	{
		Name:     "unsignedLong",
		BaseType: "nonNegativeInt",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"maxInclusive", "18446744073709551615"},
			{"minInclusive", "0"},
		},
	},
	{
		Name:     "unsignedInt",
		BaseType: "unsignedLong",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"maxInclusive", "4294967295"},
			{"minInclusive", "0"},
		},
	},
	{
		Name:     "unsignedShort",
		BaseType: "unsignedInt",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"maxInclusive", "65535"},
			{"minInclusive", "0"},
		},
	},
	{
		Name:     "unsignedByte",
		BaseType: "unsignedShort",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"maxInclusive", "255"},
			{"minInclusive", "0"},
		},
	},
	{
		Name:     "positiveInteger",
		BaseType: "nonNegativeInt",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"minInclusive", "1"},
		},
	},
	{
		Name:     "yearMonthDur",
		BaseType: "duration",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"pattern", `[^DT]*`},
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "dayTimeDur",
		BaseType: "duration",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"pattern", `[^YM]*(T.*)?`},
			{"whiteSpace", "collapse"},
		},
	},
	{
		Name:     "dateTimeStamp",
		BaseType: "dateTime",
		Variety:  VarietyAtomic,
		Facets: []Facet{
			{"explicitTimezone", "required"},
		},
	},
}
