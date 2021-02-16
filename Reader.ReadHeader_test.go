package mxt

import (
	"strings"
	"testing"
)

func TestValidHeaders(t *testing.T) {
	var allTests = []struct {
		input   string
		name    string
		comment string
	}{
		{
			`// filename -->`,
			"filename",
			"",
		},
		{
			`// file.name this is a comment -->`,
			"file.name",
			"this is a comment",
		},
		{
			`// file/name
// this is a comment line -->`,
			"file/name",
			"this is a comment line",
		},
		{
			`// file.name this is a comment line 1
// and line 2 -->`,
			"file.name",
			"this is a comment line 1 and line 2",
		},
		{
			`// file-name.txt this is a comment line 1
// and line 2
//
// line 3
// -->`,
			"file-name.txt",
			"this is a comment line 1 and line 2\nline 3",
		},
		{
			`//---------------------     file.name.txt      
// ------------------------------------------------>`,
			"file.name.txt",
			"",
		},
		{
			`// file@sys:cfg --> XYZ`,
			"file@sys:cfg",
			"",
		},
	}

	for _, test := range allTests {
		r := NewReader(strings.NewReader(test.input))
		c, err := r.ReadChunk()
		h := c.Header
		if h.Name != test.name {
			t.Errorf("invalid name: %q != %q", h.Name, test.name)
		}
		if h.Comment != test.comment {
			t.Errorf("invalid comment: %q != %q", h.Comment, test.comment)
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestInvalidHeader(t *testing.T) {
	var allTests = []struct {
		input string
		err   error
	}{}

	for _, test := range allTests {
		r := NewReader(strings.NewReader(test.input))
		_, err := r.ReadChunk()
		if err != test.err {
			t.Errorf("invalid error: %v != %v", err, test.err)
		}
	}
}
