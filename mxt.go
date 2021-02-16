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

// A Chunk represents a mxt chunk.
type Chunk struct {
	Header
	Content string
}
