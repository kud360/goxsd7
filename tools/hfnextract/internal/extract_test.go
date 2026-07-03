package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractTables(t *testing.T) {
	content := `
# Header
Some text.

| Col1 | Col2 |
| --- | --- |
| Val1 | Val2 |
| Val3 | Val4 |

Other text.

| ColA | ColB |
| --- | --- |
| VA | VB |
`
	tmpDir := t.TempDir()
	fpath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		headerText string
		index      int
		want       []Table
	}{
		{
			name:       "extract all",
			index:      -1,
			headerText: "",
			want: []Table{
				{
					Header: []string{"Col1", "Col2"},
					Rows: [][]string{
						{"Val1", "Val2"},
						{"Val3", "Val4"},
					},
				},
				{
					Header: []string{"ColA", "ColB"},
					Rows: [][]string{
						{"VA", "VB"},
					},
				},
			},
		},
		{
			name:       "extract by index 0",
			index:      0,
			headerText: "",
			want: []Table{
				{
					Header: []string{"Col1", "Col2"},
					Rows: [][]string{
						{"Val1", "Val2"},
						{"Val3", "Val4"},
					},
				},
			},
		},
		{
			name:       "extract by index 1",
			index:      1,
			headerText: "",
			want: []Table{
				{
					Header: []string{"ColA", "ColB"},
					Rows: [][]string{
						{"VA", "VB"},
					},
				},
			},
		},
		{
			name:       "extract by header partial match",
			index:      -1,
			headerText: "ColA",
			want: []Table{
				{
					Header: []string{"ColA", "ColB"},
					Rows: [][]string{
						{"VA", "VB"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTables(fpath, tt.headerText, tt.index)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
