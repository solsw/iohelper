package iohelper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"
	"sync/atomic"

	"github.com/solsw/errorhelper"
)

type lastResult struct {
	n   int
	err error
}

// NonBlockWriter is a wrapper around [io.Writer] that does not block on [io.Writer.Write] call.
// NonBlockWriter implements [io.WriteCloser] interface.
// To gracefully wait for the wrapped [io.Writer] to finish writing,
// [NonBlockWriter.Close] method must be called (typically by [defer] statement).
//
// [defer]: https://go.dev/ref/spec#Defer_statements
type NonBlockWriter struct {
	ctx       context.Context
	wr        io.Writer
	ch        chan []byte
	lr        atomic.Pointer[lastResult]
	done      sync.WaitGroup
	closed    atomic.Bool
	onceClose sync.Once
	writing   atomic.Bool
}

// NewNonBlockWriter creates a new [NonBlockWriter].
// 'w' - wrapped [io.Writer].
// 'size' - buffer size of the underlying chan []byte (defaults to [math.MaxInt16] if zero or negative).
// 'onError' (if not nil) is called if the wrapped [io.Writer.Write] returns an error.
// If 'onError' returns true (or ctx.Err() returns a non-nil error), [NonBlockWriter] is closed
// and any remaining non-written data is discarded by [NonBlockWriter.Close] method.
func NewNonBlockWriter(ctx context.Context, w io.Writer, size int, onError func(error) bool) (*NonBlockWriter, error) {
	if w == nil {
		return nil, errorhelper.CallerError(errors.New("nil writer"))
	}
	if size < 1 {
		size = math.MaxInt16
	}
	newnbw := NonBlockWriter{
		ctx: ctx,
		wr:  w,
		ch:  make(chan []byte, size),
	}
	newnbw.done.Add(1)
	go func(nbw *NonBlockWriter) {
		for bb := range nbw.ch {
			if ctxErr := nbw.ctx.Err(); ctxErr != nil {
				nbw.lr.Store(&lastResult{err: ctxErr})
				break
			}
			nbw.writing.Store(true)
			n, err := nbw.wr.Write(bb)
			nbw.lr.Store(&lastResult{n: n, err: err})
			nbw.writing.Store(false)
			if err != nil && onError != nil && onError(err) {
				break
			}
		}
		nbw.closed.Store(true)
		nbw.done.Done()
	}(&newnbw)
	return &newnbw, nil
}

// Write implements the [io.Writer] interface.
// Write returns an error only if 'nbw' is closed.
// An error (if any) returned by the wrapped [io.Writer.Write] call is returned
// by [NonBlockWriter.Close] or [NonBlockWriter.LastResult] methods.
func (nbw *NonBlockWriter) Write(p []byte) (int, error) {
	if nbw.closed.Load() {
		if lr := nbw.lr.Load(); lr != nil && lr.err != nil {
			return 0, fmt.Errorf("NonBlockWriter is closed: %w", lr.err)
		}
		return 0, errors.New("NonBlockWriter is closed")
	}
	if len(p) == 0 {
		return 0, nil
	}
	// since the same slice may be passed to this method in separate calls (e.g. as [log.Println] does),
	// the current contents of 'p' must be copied to a new local slice
	// fmt.Printf("%p\n", p) <- prints the same address when NonBlockWriter is passed to log.SetOutput
	locp := make([]byte, len(p))
	// fmt.Printf("%p\n", locp) <- prints different addresses
	copy(locp, p)
	nbw.ch <- locp
	return len(p), nil
}

// Close waits for the wrapped [io.Writer] to finish writing
// and returns the error (if any) returned by the last wrapped [io.Writer.Write] call.
// Close implements the [io.Closer] interface.
func (nbw *NonBlockWriter) Close() error {
	nbw.onceClose.Do(func() {
		nbw.closed.Store(true)
		close(nbw.ch)
		nbw.done.Wait()
	})
	if lr := nbw.lr.Load(); lr != nil {
		return lr.err
	}
	return nil
}

// LastResult returns result of the last wrapped [io.Writer.Write] call.
func (nbw *NonBlockWriter) LastResult() (int, error) {
	if lr := nbw.lr.Load(); lr != nil {
		return lr.n, lr.err
	}
	return 0, nil
}

// IsWriting reports whether the wrapped [io.Writer] is in the writing phase or not.
func (nbw *NonBlockWriter) IsWriting() bool {
	return nbw.writing.Load()
}
