// Package mxt reads and writes mxt files.
//
// See https://mxt.aiq.dk// for more information about the mxt file format.
package mxt

//********************************************************************** Header

// A Header represents the user information in a mxt header.
type Header struct {
	Name    string
	Comment string
}

func (h Header) IsEmpty() bool {
	return h.Name == "" && h.Comment == ""
}

// A Chunk represents a mxt chunk.
type Chunk struct {
	Name    string
	Comment string
	Content string
}

func (c Chunk) Header() Header {
	return Header{
		Name:    c.Name,
		Comment: c.Comment,
	}
}

func (c Chunk) IsEmpty() bool {
	return c.Name == "" && c.Comment == "" && c.Content == ""
}

type Chunks []Chunk

func (cs Chunks) Get(name string) (Chunk, bool) {
	for _, c := range cs {
		if c.Name == name {
			return c, true
		}
	}
	return Chunk{}, false
}

func (cs Chunks) Find(name string) int {
	for i, c := range cs {
		if c.Name == name {
			return i
		}
	}
	return -1
}

func (cs Chunks) Set(name string, content string) Chunks {
	nc := Chunk{
		Name:    name,
		Content: content,
	}
	for i, c := range cs {
		if c.Name == name {
			cs[i] = nc
		}
	}
	cs = append(cs, nc)
	return cs
}
