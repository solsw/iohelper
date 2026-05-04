package iohelper

import (
	"errors"
	"io"

	"github.com/solsw/errorhelper"
)

// NilSafeWriter is an [io.Writer] implementation,
// whose Write method does nothing if underlying [io.Writer] is nil.
type NilSafeWriter struct {
	writer io.Writer
}

// NewNilSafeWriter creates a new [NilSafeWriter] based on 'w'.
func NewNilSafeWriter(w io.Writer) *NilSafeWriter {
	return &NilSafeWriter{writer: w}
}

// Write implements the [io.Writer] interface.
func (nw *NilSafeWriter) Write(p []byte) (int, error) {
	if nw == nil {
		return 0, errorhelper.CallerError(errors.New("nil NilSafeWriter"))
	}
	if nw.writer == nil {
		return len(p), nil
	}
	return nw.writer.Write(p)
}
