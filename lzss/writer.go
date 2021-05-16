package lzss

import (
	"bufio"
	"io"
)

type writer interface {
	io.ByteWriter
	Flush() error
}

type compressor struct {
	w         writer
	windowLen int
}

func (c *compressor) Write(p []byte) (n int, err error) {
	for _, b := range p {
		err = c.w.WriteByte(b)
		if err != nil {
			return n, err
		}
		n += 1
	}

	return n, err
}

func (c *compressor) Close() error {
	return c.w.Flush()
}

func NewWriter(w io.Writer) io.WriteCloser {
	return NewWriterWindow(w, 4096)
}

func NewWriterWindow(w io.Writer, windowLen int) io.WriteCloser {

	return &compressor{
		w:         bufio.NewWriter(w),
		windowLen: windowLen,
	}
}
