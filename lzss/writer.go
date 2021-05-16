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
	 w writer
	 windowLen int
	 err error
}

func (c *compressor) Write(p []byte) (n int, err error) {
	for _,b := range p {
		err = c.w.WriteByte(b)
		if err != nil {
			n += 1
		}
	}

	return n, err
}

func (c *compressor) Close() error {
	return c.w.Flush()
}

func NewWriter(w io.Writer, windowLen int) io.WriteCloser {
	bw, ok := w.(writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}

	return &compressor{
		w: bw,
		windowLen: windowLen,
	}
}
