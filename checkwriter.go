package iohelper

import (
	"errors"
	"io"

	"github.com/solsw/errorhelper"
)

// CheckWriter is an [io.Writer] implementation, that checks
// the underlying [io.Writer] and slice data by calling
// the provided 'check' function before each [io.Writer.Write] call.
type CheckWriter struct {
	writer io.Writer
	check  func(*io.Writer, []byte) error
}

// NewCheckWriter creates a new [CheckWriter] based on 'w' and 'check'.
// If 'check' returns an error, this error is returned by [Write] method.
// 'check' may replace the provided [io.Writer].
func NewCheckWriter(w io.Writer, check func(*io.Writer, []byte) error) (*CheckWriter, error) {
	if w == nil {
		return nil, errorhelper.CallerError(errors.New("nil writer"))
	}
	if check == nil {
		return nil, errorhelper.CallerError(errors.New("nil check"))
	}
	return &CheckWriter{writer: w, check: check}, nil
}

// Write implements the [io.Writer] interface.
func (cw *CheckWriter) Write(p []byte) (int, error) {
	if cw == nil {
		return 0, errorhelper.CallerError(errors.New("nil CheckWriter"))
	}
	if cw.check == nil {
		return 0, errorhelper.CallerError(errors.New("nil check"))
	}
	if err := cw.check(&cw.writer, p); err != nil {
		return 0, errorhelper.CallerError(err)
	}
	if cw.writer == nil {
		return 0, errorhelper.CallerError(errors.New("nil writer"))
	}
	return cw.writer.Write(p)
}
