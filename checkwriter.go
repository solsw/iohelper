package iohelper

import (
	"io"
)

// CheckWriter is an [io.Writer] implementation,
// that checks the underlying [io.Writer] before each [io.Writer.Write] call.
type CheckWriter struct {
	writer io.Writer
	check  func(*io.Writer) error
}

// NewCheckWriter creates a new [CheckWriter] based on 'w' and 'check'.
// If 'check' returns an error, this error is returned by [Write] method.
func NewCheckWriter(w io.Writer, check func(*io.Writer) error) *CheckWriter {
	return &CheckWriter{writer: w, check: check}
}

// Write implements the [io.Writer] interface.
func (cw *CheckWriter) Write(p []byte) (int, error) {
	if err := cw.check(&cw.writer); err != nil {
		return 0, err
	}
	return cw.writer.Write(p)
}
