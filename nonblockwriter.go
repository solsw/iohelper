package iohelper

import (
	"errors"
	"io"
	"math"
	"sync"
)

// NonBlockWriter is a wrapper around io.Writer that does not block on Write call.
// To gracefully wait for the wrapped io.Writer to finish writing, Close must be called (typically by 'defer' statement).
type NonBlockWriter struct {
	wr        io.Writer
	ch        chan []byte
	n         int
	err       error
	done      sync.WaitGroup
	closed    bool
	onceClose sync.Once
	writing   bool
}

// NewNonBlockWriter creates a new NonBlockWriter.
// 'w' - wrapped io.Writer.
// 'size' - buffer size of the underlying chan []byte (defaults to math.MaxInt16 if zero or negative).
// 'onError' (if not nil) is called if the wrapped io.Writer.Write returns an error.
// If 'onError' returns true, NonBlockWriter is immediately closed and any remaining non-written data is discarded.
func NewNonBlockWriter(w io.Writer, size int, onError func(error) bool) *NonBlockWriter {
	if size < 1 {
		size = math.MaxInt16
	}
	newnbw := NonBlockWriter{
		wr: w,
		ch: make(chan []byte, size),
	}
	newnbw.done.Add(1)
	go func(nbw *NonBlockWriter) {
		for bb := range nbw.ch {
			nbw.writing = true
			nbw.n, nbw.err = nbw.wr.Write(bb)
			nbw.writing = false
			if nbw.err != nil && onError != nil && onError(nbw.err) {
				nbw.done.Done()
				nbw.Close()
				return
			}
		}
		nbw.done.Done()
	}(&newnbw)
	return &newnbw
}

// Write implements the io.Writer interface.
// If NonBlockWriter is not closed, Write returns len(p) and nil.
func (nbw *NonBlockWriter) Write(p []byte) (int, error) {
	if nbw.closed {
		return -1, errors.New("NonBlockWriter is closed")
	}
	if len(p) == 0 {
		return 0, nil
	}
	// since the same slice may be passed to this method in separate calls (e.g. as log.Println does),
	// the current contents of 'p' must be copied to a new local slice
	// fmt.Printf("%p\n", p) <- prints the same address when NonBlockWriter is passed to log.SetOutput
	locp := make([]byte, len(p))
	// fmt.Printf("%p\n", locp) <- prints different addresses
	copy(locp, p)
	nbw.ch <- locp
	return len(p), nil
}

// Close waits for the wrapped io.Writer to finish writing
// and returns the error (if any) returned by the last wrapped io.Writer.Write call.
// Close implements the io.Closer interface.
func (nbw *NonBlockWriter) Close() error {
	nbw.onceClose.Do(func() {
		nbw.closed = true
		close(nbw.ch)
		nbw.done.Wait()
	})
	return nbw.err
}

// LastResult returns result of the last wrapped io.Writer.Write call.
func (nbw *NonBlockWriter) LastResult() (int, error) {
	return nbw.n, nbw.err
}

// IsWriting reports whether the wrapped io.Writer is in the writing phase or not.
func (nbw *NonBlockWriter) IsWriting() bool {
	return nbw.writing
}
