package xmltree

import (
	"encoding/xml"
	"io"
	"sort"

	"github.com/kud360/goxsd7/xsderr"
)

// Token wraps an encoding/xml.Token with its source location.
type Token struct {
	token xml.Token
	loc   xsderr.Loc
}

func (t Token) Value() xml.Token { return t.token }
func (t Token) Loc() xsderr.Loc  { return t.loc }

// Decoder is a streaming XML reader that tracks line and column positions.
type Decoder struct {
	dec *xml.Decoder
	pt  *positionTracker
	uri string
}

// NewDecoder creates a new Decoder for the given reader and URI.
func NewDecoder(r io.Reader, uri string) *Decoder {
	pt := &positionTracker{
		r: r,
	}
	return &Decoder{
		dec: xml.NewDecoder(pt),
		pt:  pt,
		uri: uri,
	}
}

// Decode returns the next XML token and its location.
func (d *Decoder) Decode() (Token, error) {
	off := d.dec.InputOffset()
	t, err := d.dec.Token()
	if err != nil {
		return Token{}, err
	}

	return Token{
		token: t,
		loc:   d.pt.offsetToLoc(off, d.uri),
	}, nil
}

type positionTracker struct {
	r           io.Reader
	newlineOffs []int64
	currentOff  int64
}

func (pt *positionTracker) Read(p []byte) (n int, err error) {
	n, err = pt.r.Read(p)
	for i := 0; i < n; i++ {
		if p[i] == '\n' {
			pt.newlineOffs = append(pt.newlineOffs, pt.currentOff+int64(i))
		}
	}
	pt.currentOff += int64(n)
	return n, err
}

func (pt *positionTracker) offsetToLoc(off int64, uri string) xsderr.Loc {
	idx := sort.Search(len(pt.newlineOffs), func(i int) bool {
		return pt.newlineOffs[i] >= off
	})

	line := idx + 1
	col := int(off) + 1
	if idx > 0 {
		lastNewlineOff := pt.newlineOffs[idx-1]
		col = int(off - lastNewlineOff)
	}

	return xsderr.Loc{
		URI:  uri,
		Line: line,
		Col:  col,
	}
}
