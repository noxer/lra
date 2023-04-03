package lra

import "io"

type LazyReaderAt struct {
	r   io.Reader
	buf []byte
}

func NewLazyReaderAt(r io.Reader) *LazyReaderAt {
	return &LazyReaderAt{r: r}
}

func (l *LazyReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	// calculate how many bytes we need
	expect := int64(len(p)) + off
	diff := expect - int64(len(l.buf))
	if diff > 0 {
		// we need to read more bytes
		err = l.readMore(diff)

		// we didn't manage to read any bytes of the ones the caller requested
		if int64(len(l.buf)) <= off {
			return 0, err
		}
	}

	// copy the bytes into the readerAt's buffer. copy will only read bytes until callers buffer is full
	// or our buffer has no more bytes and returns the number of bytes copied. exactly what we need for n.
	n = copy(p, l.buf[off:])
	return
}

func (l *LazyReaderAt) readMore(n int64) error {
	// expand the buffer to hold our new data
	newBuf := append(l.buf, make([]byte, n)...)

	// try to read the data into the new buffer
	nn, err := io.ReadFull(l.r, newBuf[len(l.buf):])

	// update the buffer with the new length
	l.buf = newBuf[:len(l.buf)+nn]

	return err
}

// Reset prepares the LazyReaderAt for use with a new io.Reader.
func (l *LazyReaderAt) Reset(r io.Reader) {
	// trim the buffer without freeing it
	l.buf = l.buf[:0]
	l.r = r
}

// make sure our struct implements the interface
var _ io.ReaderAt = (*LazyReaderAt)(nil)
