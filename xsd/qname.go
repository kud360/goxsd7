package xsd

// Namespace URIs fixed by the specifications.
const (
	XSDNamespace = "http://www.w3.org/2001/XMLSchema"
	XSINamespace = "http://www.w3.org/2001/XMLSchema-instance"
	XMLNamespace = "http://www.w3.org/XML/1998/namespace"
)

// QName is an expanded name: namespace URI plus local part.
// The zero QName is the absent name (unnamed/anonymous components).
type QName struct {
	Space string // namespace URI; "" = no namespace
	Local string
}

// Builtin returns the QName of a builtin datatype by local name.
func Builtin(local string) QName {
	return QName{Space: XSDNamespace, Local: local}
}

func (q QName) IsZero() bool { return q == QName{} }

// String renders Clark notation: {namespace}local.
func (q QName) String() string {
	if q.Space == "" {
		return q.Local
	}
	return "{" + q.Space + "}" + q.Local
}
