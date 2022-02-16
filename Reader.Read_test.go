package mxt

import (
	"io"
	"strings"
	"testing"
)

type TestToken struct {
	res map[string]string
	err error
	*testing.T
}

func (tt *TestToken) expectErr(err error) {
	if tt.err != err {
		tt.Errorf("invalid error: %v != %v", tt.err, err)
	}
}

func (tt *TestToken) expectLen(n int) {
	if len(tt.res) != n {
		tt.Errorf("invalid map lengt: %q != %q", len(tt.res), n)
	}
}

func (tt *TestToken) expectContent(name string, content string) {
	cnt, found := tt.res[name]
	if !found {
		tt.Errorf("no entry for: %q", name)
	}
	if cnt != content {
		tt.Errorf("invalid content: \n%s\n%s", cnt, content)
	}
}

//******************************************************************************

func TestReadExample(t *testing.T) {
	input := `//---------------------------------------------------------------- user.json -->
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
}`

	//**************************************************************************

	res, err := Read(strings.NewReader(input))

	tt := &TestToken{res, err, t}
	tt.expectLen(5)
	tt.expectErr(io.EOF)

	tt.expectContent("user.json", `{
   "user": "alucard",
   "password": "C:SotN1997"
}`)

	tt.expectContent("connection.ini", `request: GET

[url]
schema=http
host=localhost
port=8080
path=/db/add`)

	tt.expectContent("user.pgp", `-----BEGIN PGP MESSAGE-----

hQEMA8p144+Gi+YpAQf/VeFG9Zb+8w9aldWll8n2g3jqpE613LKg2XAJgwXQmSQL
R4O+TlQakJ+Mz5vM4IxxubPgYCyt6cyL7qM3oJIuk7vsqMbl5t7c/dOfXjj7goIC
IskIX+9e5qrr8jRG/KZYSdBJtFEI9oNtZTLlnv3yeV3OWNTbUnjdTWrk/h1kavJE
D8nD4suo6ckVzYGJpknGSIAwaCFl//aqR/3SWO4wi6ibbfub8LA73V90Ll3/S/Ph
xU15HYmdCATnVX1sp1PWmyz972bMvl8txyIKMUueVw+w0C19ZTfWXjuFSguF7zt7
RY+I3to2lbyVJbcI9Dyz04GOJZ2vIhG9eq65FxeweAKDa7L+iH1NA5L2lYd9DEr1
ro/CU6vIqkOSNRUrNYDwqz1g3Z3eAQB/8t9Y4WsV4KL0M229rsFrtl26i7+quYfg
uuTd
=WxK9
-----END PGP MESSAGE-----`)

	tt.expectContent("hello-world.h", "")

	tt.expectContent("hello-world.c", `// this is part of hello-world.c
#include<stdio.h>

int main(void) {
   printf("Hello World\n");
   return 0;
}`)
}
