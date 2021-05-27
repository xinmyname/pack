package lzss

import (
	"io"

	"github.com/xinmyname/bitstream-go"
)

func NewReader(r io.Reader) io.ReadCloser {
	const maxWindowSize = 255

	return &decompressor{
		bs:         *bitstream.NewReader(r, &bitstream.ReaderOptions{BufferSize: 256}),
		windowSize: maxWindowSize,
		buffer:     make([]uint8, 0, maxWindowSize),
	}
}

type decompressor struct {
	bs         bitstream.Reader
	windowSize int
	buffer     []uint8
	pos        int
	wrote      int
}

func (d *decompressor) Read(p []byte) (int, error) {

	d.wrote = 0

	for d.wrote < len(p) {

		tokenType, err := d.bs.ReadBit()

		if err != nil {
			return d.wrote, err
		}

		var length uint8 = 0

		if tokenType&0x1 == 0x1 {
			offset, err := d.bs.ReadUint8()

			if err != nil {
				return d.wrote, err
			}

			length, err = d.bs.ReadUint8()

			if err != nil {
				return d.wrote, err
			}

			if length == 0 && offset == 0 {
				return d.wrote, io.EOF
			}

			offset = uint8(len(d.buffer) - int(offset))

			d.buffer = append(d.buffer, d.buffer[offset:offset+length]...)
		} else {
			byte, err := d.bs.ReadUint8()

			if err != nil {
				return d.wrote, err
			}
			length = 1
			d.buffer = append(d.buffer, byte)
		}

		if len(d.buffer) > d.windowSize {
			newstart := len(d.buffer) - d.windowSize
			d.buffer = d.buffer[newstart:]
			d.pos -= newstart
		}

		advance := copy(p[d.wrote:], d.buffer[d.pos:])
		d.pos += advance
		d.wrote += advance
	}

	return d.wrote, nil
}

func (d *decompressor) Close() error {
	return nil
}
