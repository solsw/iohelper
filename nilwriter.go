package iohelper

import (
	"io"
)

// NilWriter is an [io.Writer] implementation,
// whose Write method does nothing if underlying [io.Writer] is nil.
type NilWriter struct {
	writer io.Writer
}

// NewNilWriter creates a new [NilWriter] based on 'w'.
func NewNilWriter(w io.Writer) *NilWriter {
	return &NilWriter{writer: w}
}

// Write implements the [io.Writer] interface.
func (nw *NilWriter) Write(p []byte) (int, error) {
	if nw.writer == nil {
		return len(p), nil
	}
	return nw.writer.Write(p)
}
