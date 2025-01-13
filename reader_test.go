package mxt

import (
	"strings"
	"testing"
)

func TestReadExample(t *testing.T) {
	cases := []struct {
		inp string
		exp Chunks
	}{
		{
			inp: `//---------------------------------------------------------------- user.json -->
{
   "user": "alucard",
   "password": "C:SotN1997"
}
//--------------------------------------------------------------- connection.ini
// comment line that is not part of the ini file,
// comment lines will be joined with a space character
//
// empty comment lines will generate a newline character in the comment
//----------------------------------------------------------------------------->
request: GET

[url]
schema=http
host=localhost
port=8080
path=/db/add
// user.pgp --> XYZ
-----BEGIN PGP MESSAGE-----

hQEMA8p144+Gi+YpAQf/VeFG9Zb+8w9aldWll8n2g3jqpE613LKg2XAJgwXQmSQL
R4O+TlQakJ+Mz5vM4IxxubPgYCyt6cyL7qM3oJIuk7vsqMbl5t7c/dOfXjj7goIC
IskIX+9e5qrr8jRG/KZYSdBJtFEI9oNtZTLlnv3yeV3OWNTbUnjdTWrk/h1kavJE
D8nD4suo6ckVzYGJpknGSIAwaCFl//aqR/3SWO4wi6ibbfub8LA73V90Ll3/S/Ph
xU15HYmdCATnVX1sp1PWmyz972bMvl8txyIKMUueVw+w0C19ZTfWXjuFSguF7zt7
RY+I3to2lbyVJbcI9Dyz04GOJZ2vIhG9eq65FxeweAKDa7L+iH1NA5L2lYd9DEr1
ro/CU6vIqkOSNRUrNYDwqz1g3Z3eAQB/8t9Y4WsV4KL0M229rsFrtl26i7+quYfg
uuTd
=WxK9
-----END PGP MESSAGE-----
//XYZ hello-world.h -->
//---------------------------------------------------------- hello-world.c --> X
// this is part of hello-world.c
#include<stdio.h>

int main(void) {
   printf("Hello World\n");
   return 0;
}`,
			exp: Chunks{
				{
					Name: "user.json",
					Content: `{
   "user": "alucard",
   "password": "C:SotN1997"
}`,
				},
				{
					Name: "connection.ini",
					Content: `request: GET

[url]
schema=http
host=localhost
port=8080
path=/db/add`,
				},
				{
					Name: "user.pgp",
					Content: `-----BEGIN PGP MESSAGE-----

hQEMA8p144+Gi+YpAQf/VeFG9Zb+8w9aldWll8n2g3jqpE613LKg2XAJgwXQmSQL
R4O+TlQakJ+Mz5vM4IxxubPgYCyt6cyL7qM3oJIuk7vsqMbl5t7c/dOfXjj7goIC
IskIX+9e5qrr8jRG/KZYSdBJtFEI9oNtZTLlnv3yeV3OWNTbUnjdTWrk/h1kavJE
D8nD4suo6ckVzYGJpknGSIAwaCFl//aqR/3SWO4wi6ibbfub8LA73V90Ll3/S/Ph
xU15HYmdCATnVX1sp1PWmyz972bMvl8txyIKMUueVw+w0C19ZTfWXjuFSguF7zt7
RY+I3to2lbyVJbcI9Dyz04GOJZ2vIhG9eq65FxeweAKDa7L+iH1NA5L2lYd9DEr1
ro/CU6vIqkOSNRUrNYDwqz1g3Z3eAQB/8t9Y4WsV4KL0M229rsFrtl26i7+quYfg
uuTd
=WxK9
-----END PGP MESSAGE-----`,
				},
				{
					Name:    "hello-world.h",
					Content: "",
				},
				{
					Name: "hello-world.c",
					Content: `// this is part of hello-world.c
#include<stdio.h>

int main(void) {
   printf("Hello World\n");
   return 0;
}`,
				},
			},
		},
	}

	for i, c := range cases {
		chunks, err := ReadString(c.inp)
		if err != nil {
			t.Errorf("%v - unexpected error: %v", i, err)
			continue
		}
		if len(chunks) != len(c.exp) {
			t.Errorf("%v - unexpected number chunks: %v != %v", i, len(chunks), len(c.exp))
			continue
		}

		for ci, cv := range chunks {
			expName := c.exp[ci].Name
			if expName != cv.Name {
				t.Errorf("%v - @%v - unexpected name: %q != %q", i, ci, expName, cv.Name)
				break
			}
			expContent := c.exp[ci].Content
			if expContent != cv.Content {
				t.Errorf("%v - @%v - unexpected content: %q != %q", i, ci, expContent, cv.Content)
				break
			}
		}
	}
}

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
		h := c.Header()
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
