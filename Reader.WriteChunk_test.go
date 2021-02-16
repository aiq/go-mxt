package mxt

import (
	"strings"
	"testing"
)

func TestValidChunks(t *testing.T) {
	var allTests = []struct {
		name    string
		comment string
		content string
		exp     string
	}{
		{
			"filename",
			"",
			"bla bla bla",
			`// filename -->
bla bla bla`,
		},
		{
			"the.name",
			"a comment\nwith multiple lines",
			"",
			`// the.name
// a comment
//
// with multiple lines
// -->
`,
		},
		{
			"a.b.c",
			"a long comment that will be split between multiple lines to be not longer as 80 characters in one line",
			"// some content \nwith backslash at the begin",
			`// a.b.c
// a long comment that will be split between multiple lines to be not longer as
// 80 characters in one line
// --> XYZ
// some content 
with backslash at the begin`,
		},
	}

	for _, test := range allTests {
		var b strings.Builder
		w := NewWriter(&b)

		_, err := w.WriteChunk(Chunk{Header{test.name, test.comment}, test.content})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		res := b.String()
		if res != test.exp {
			t.Errorf("unexpected result: %s", res)
		}
	}
}
