package mxt

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

//********************************************************************** Writer

// A Writer allows to create valid mxt files.
//
// As returned by NewWriter, a Writer writes chunks with as Salt value if
// necessary.
// The Salt value will be extendet with random values from the Alphabet slice.
//
// Salt is set by NewWriter with "XYZ".
//
// Alphabet is set by NewWriter with
// "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz?!§$@€&-~".
//
// If an error occurs writing to a Writer, no more chunks will be accepted and
// all subsequent writes will return the error.
type Writer struct {
	Alphabet []rune //
	Salt     string //
	prevPatt string
	writer   *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Alphabet: []rune("AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz?!§$@€&-~"),
		Salt:     "XYZ",
		prevPatt: "",
		writer:   bufio.NewWriter(w),
	}
}

//***************************************************************** Writer.func

func (my *Writer) detectPattern(content string) string {
	if !strings.HasPrefix(content, "//") && !strings.Contains(content, "\n//") {
		return ""
	}

	patt := []rune(my.Salt)
	for {
		pattStr := string(patt)
		if !strings.Contains(content, "\n//-"+string(pattStr)+"-") {
			return pattStr
		}
		patt = append(patt, my.Alphabet[rand.Intn(len(my.Alphabet))])
	}
}

func subLines(baseLine string) []string {
	if len(baseLine) <= 77 {
		return []string{baseLine}
	}
	result := []string{}
	linelen := 0
	line := []string{}
	for _, word := range strings.Fields(baseLine) {
		n := linelen + len(line)

		if n+len(word) > 77 && n != 0 {
			result = append(result, strings.Join(line, " "))
			linelen = 0
			line = []string{}
		}

		linelen += len(word)
		line = append(line, word)
	}
	if linelen > 0 {
		result = append(result, strings.Join(line, " "))
	}
	return result
}

func (my *Writer) writeComment(comment string) (n int, err error) {
	if len(comment) == 0 {
		return 0, nil
	}
	tmp := 0
	lines := strings.Split(comment, "\n")
	for i, line := range lines {
		for _, li := range subLines(strings.TrimSpace(line)) {
			tmp, err = my.writer.WriteString("\n// " + li)
			n += tmp
			if err != nil {
				return n, err
			}
		}
		if i != len(lines)-1 { // last line?
			tmp, err = my.writer.WriteString("\n//")
			n += tmp
			if err != nil {
				return n, err
			}
		}
	}

	return n, nil
}

func (my *Writer) writeHeader(h Header, nextPatt string) (n int, err error) {
	if len(h.Name) == 0 {
		return 0, fmt.Errorf("empty header name")
	}
	if strings.Contains(h.Name, " ") {
		return 0, fmt.Errorf("invalid header name: %s", h.Name)
	}

	tmp := 0

	marker := "// "
	if len(my.prevPatt) > 0 {
		marker = "//-" + my.prevPatt + "- "
	}
	tmp, err = my.writer.WriteString(marker)
	n += tmp
	if err != nil {
		return n, err
	}

	tmp, err = my.writer.WriteString(h.Name)
	n += tmp
	if err != nil {
		return n, err
	}

	tmp, err = my.writeComment(h.Comment)
	n += tmp
	if err != nil {
		return n, err
	}

	my.prevPatt = nextPatt
	arrow := " -->"
	if len(nextPatt) > 0 {
		arrow = " --> " + nextPatt
	}
	if len(h.Comment) > 0 {
		arrow = "\n//" + arrow
	}
	tmp, err = my.writer.WriteString(arrow + "\n")
	n += tmp

	return n, err
}

// WriteChunk writes a single Chunk to w along with any necessary marker and
// salt values.
// Each Chunk will be written to the underlying io.Writer.
func (my *Writer) WriteChunk(c Chunk) (n int, err error) {
	tmp := 0
	tmp, err = my.writeHeader(c.Header, my.detectPattern(c.Content))
	n += tmp
	if err != nil {
		return n, err
	}
	tmp, err = my.writer.WriteString(c.Content)
	n += tmp

	err = my.writer.Flush()

	return n, err
}

// Write writes a Chunk without comment.
func (my *Writer) Write(name string, content string) (int, error) {
	return my.WriteChunk(Chunk{
		Header: Header{
			Name: name,
		},
		Content: content,
	})
}

//*********************************************************************** Write

// Write writes all map entries to w.
// The key will be used as name, the value will be stored as content in the mxt.
func Write(m map[string]string, w io.Writer) (n int, err error) {
	writer := NewWriter(w)
	tmp := 0
	for name, content := range m {
		tmp, err = writer.Write(name, content)
		n += tmp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// WriteFile writes all map entries to the file behind path.
func WriteFile(m map[string]string, path string) (n int, err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	return Write(m, file)
}
