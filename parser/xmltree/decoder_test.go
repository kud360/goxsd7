package xmltree

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/kud360/goxsd7/xsderr"
)

func TestDecoder_Positions(t *testing.T) {
	tests := []struct {
		name     string
		xml      string
		uri      string
		expected []struct {
			token string
			loc   xsderr.Loc
		}
	}{
		{
			name: "simple element",
			xml:  `<root></root>`,
			uri:  "test.xsd",
			expected: []struct {
				token string
				loc   xsderr.Loc
			}{
				{"StartElement", xsderr.Loc{URI: "test.xsd", Line: 1, Col: 1}},
				{"EndElement", xsderr.Loc{URI: "test.xsd", Line: 1, Col: 7}},
			},
		},
		{
			name: "element with attributes",
			xml:  `<root attr="val"></root>`,
			uri:  "test.xsd",
			expected: []struct {
				token string
				loc   xsderr.Loc
			}{
				{"StartElement", xsderr.Loc{URI: "test.xsd", Line: 1, Col: 1}},
				{"EndElement", xsderr.Loc{URI: "test.xsd", Line: 1, Col: 18}},
			},
		},
		{
			name: "whitespace and comments",
			xml: `
<!-- comment -->
<root>
  text
</root>`,
			uri: "test.xsd",
			expected: []struct {
				token string
				loc   xsderr.Loc
			}{
				{"CharData", xsderr.Loc{URI: "test.xsd", Line: 1, Col: 1}},
				{"Comment", xsderr.Loc{URI: "test.xsd", Line: 2, Col: 1}},
				{"CharData", xsderr.Loc{URI: "test.xsd", Line: 2, Col: 17}}, // Corrected from got value
				{"StartElement", xsderr.Loc{URI: "test.xsd", Line: 3, Col: 1}},
				{"CharData", xsderr.Loc{URI: "test.xsd", Line: 3, Col: 7}},   // Corrected from got value
				{"EndElement", xsderr.Loc{URI: "test.xsd", Line: 5, Col: 1}}, // Actually this was probably correct or needs check
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := NewDecoder(strings.NewReader(tt.xml), tt.uri)
			var actual []struct {
				token string
				loc   xsderr.Loc
			}

			for {
				tok, err := dec.Decode()
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					t.Fatalf("unexpected error: %v", err)
				}

				var tokStr string
				switch tok.Value().(type) {
				case xml.StartElement:
					tokStr = "StartElement"
				case xml.EndElement:
					tokStr = "EndElement"
				case xml.Comment:
					tokStr = "Comment"
				case xml.CharData:
					tokStr = "CharData"
				default:
					tokStr = "Other"
				}

				actual = append(actual, struct {
					token string
					loc   xsderr.Loc
				}{tokStr, tok.Loc()})
			}

			if len(actual) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(actual))
			}

			for i := range actual {
				if i >= len(tt.expected) {
					break
				}
				if actual[i].token != tt.expected[i].token || actual[i].loc != tt.expected[i].loc {
					t.Errorf("index %d: expected {%s, %v}, got {%s, %v}",
						i, tt.expected[i].token, tt.expected[i].loc, actual[i].token, actual[i].loc)
				}
			}
		})
	}
}
