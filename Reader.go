package mxt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

//********************************************************************** Reader

// A Reader allows to read chunks from a mxt file.
type Reader struct {
	expPatt string
	*bufio.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		expPatt: "",
		Reader:  bufio.NewReader(r),
	}
}

//***************************************************************** Reader.func

func readString(r *bufio.Reader, delim []byte) (string, error) {
	delimByte := delim[0]
	delimTail := delim[1:]
	var builder strings.Builder
	search := true
	for search {
		tmpStr, err := r.ReadString(delimByte)
		builder.WriteString(tmpStr)
		if err != nil {
			return builder.String(), err
		}

		buf, err := r.Peek(len(delimTail))
		if err != nil {
			return builder.String(), err
		}
		if bytes.Equal(buf, delimTail) {
			search = false
			_, err = r.Read(buf)
			if err != nil {
				return builder.String(), err
			}
			_, err = builder.Write(buf)
			if err != nil {
				return builder.String(), err
			}
		}
	}
	return builder.String(), nil
}

func parseComment(rawComment string) string {
	lines := strings.Split(rawComment, "//")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
	}
	var builder strings.Builder
	space := false
	for _, line := range lines {

		s := strings.TrimSpace(line)
		if len(s) > 0 {
			if space {
				builder.WriteByte(' ')
			}
			builder.WriteString(s)
			space = true
		} else {
			builder.WriteByte('\n')
			space = false
		}
	}

	return strings.TrimSpace(builder.String())
}

func (my *Reader) readHeader() (Header, error) {
	header := Header{"", ""}
	str, err := readString(my.Reader, []byte("-->"))
	if err != nil {
		return header, err
	}

	nameBeg := strings.Index(str, " ")
	if nameBeg == -1 {
		return header, fmt.Errorf("no name beg")
	}

	str = strings.TrimSpace(str[nameBeg+1:])

	nameEnd := strings.IndexAny(str, " \n")
	if nameEnd == -1 {
		return header, fmt.Errorf("invalid header")
	}
	header.Name = str[:nameEnd]

	commentEnd := strings.LastIndexAny(str, " \n")
	if nameEnd == -1 {
		return header, fmt.Errorf("invalid header")
	}
	header.Comment = parseComment(str[nameEnd:commentEnd])

	my.expPatt, err = my.ReadString('\n')
	if err != nil && err != io.EOF {
		return header, err
	}
	my.expPatt = strings.TrimSpace(my.expPatt)

	return header, nil
}

func trimLastNewLineSuffix(str string) string {
	if strings.HasSuffix(str, "\r\n") {
		return strings.TrimSuffix(str, "\r\n")
	}
	return strings.TrimSuffix(str, "\n")
}

func (my *Reader) readContent() (string, error) {
	delim := "//"
	if my.expPatt != "" {
		delim = "//" + my.expPatt
	}

	cnt, err := readString(my.Reader, []byte(delim))
	if err != nil && err != io.EOF {
		return "", err
	}
	return trimLastNewLineSuffix(strings.TrimSuffix(cnt, delim)), nil
}

// ReadChunk reads one Chunk from r.
// If there is no data left to be read, ReadChunk returns a empty Chunk, io.EOF.
func (my *Reader) ReadChunk() (ch Chunk, err error) {

	ch.Header, err = my.readHeader()
	if err == nil {
		ch.Content, err = my.readContent()
	}

	return ch, err
}

//************************************************************************ Read

// Read creates a Reader and reads all chunks from r and stores them in the map.
// The name will be used as key, the content will be stored as value.
// Possible comments in the mxt will be ignored.
//
// A successful call returns err == nil, not err == io.EOF.
// Because Read is defined to read until EOF, it does not treat end of file as
// an error to be reported.
func Read(r io.Reader) (map[string]string, error) {
	res := make(map[string]string)
	var err error

	my := NewReader(r)
	for err == nil {
		var chunk Chunk
		chunk, err = my.ReadChunk()
		if err == nil || err != io.EOF {
			res[chunk.Header.Name] = chunk.Content
		}
	}

	return res, err
}

// ReadString reads str with Read.
func ReadString(str string) (map[string]string, error) {
	return Read(strings.NewReader(str))
}

// ReadFile reads the mxt file behind path with Read.
func ReadFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Read(file)
}
