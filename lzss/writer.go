package lzss

import (
	"io"

	"github.com/xinmyname/bitstream-go"
)

func NewWriter(w io.Writer) io.WriteCloser {
	const maxWindowSize = 255

	return &compressor{
		bs:              *bitstream.NewWriter(w),
		windowSize:      maxWindowSize,
		searchBuffer:    make([]int, 0, maxWindowSize),
		checkCharacters: make([]int, 0, maxWindowSize),
	}
}

type compressor struct {
	bs              bitstream.Writer
	windowSize      int
	searchBuffer    []int
	checkCharacters []int
}

func (c *compressor) Write(p []byte) (n int, err error) {

	compressed := c.compress(p)

	for _, token := range compressed {

		if token > 0xff {
			c.bs.WriteBit(1)
			c.bs.WriteUint8(uint8(token >> 8))   // offset
			c.bs.WriteUint8(uint8(token & 0xff)) // length
		} else {
			c.bs.WriteBit(0)
			c.bs.WriteUint8(uint8(token))
		}
	}

	// EOF sentinel
	c.bs.WriteBit(1)
	c.bs.WriteUint8(0)
	c.bs.WriteUint8(0)

	return len(p), err
}

func (c *compressor) Close() error {
	return c.bs.Flush()
}

func (c *compressor) compress(p []byte) []int {

	output := make([]int, 0, c.windowSize)

	for pos, char := range p {

		index := elementsInArray(c.checkCharacters, int(char), c.searchBuffer)

		if index == -1 || pos == len(p)-1 {
			if pos == len(p)-1 && index != -1 {
				c.checkCharacters = append(c.checkCharacters, int(char))
			}

			if len(c.checkCharacters) > 1 {
				index = elementsInArray(c.checkCharacters, -1, c.searchBuffer)
				offset := len(c.searchBuffer) - index
				length := len(c.checkCharacters)
				token := (offset << 8) | (length & 0xff)
				output = append(output, token)
				c.searchBuffer = append(c.searchBuffer, c.checkCharacters...)
			} else {
				output = append(output, c.checkCharacters...)
				c.searchBuffer = append(c.searchBuffer, c.checkCharacters...)
			}

			c.checkCharacters = c.checkCharacters[:0]
		}

		c.checkCharacters = append(c.checkCharacters, int(char))

		if len(c.searchBuffer) >= c.windowSize {
			diff := len(c.searchBuffer) - c.windowSize
			c.searchBuffer = c.searchBuffer[diff:]
		}
	}

	if len(c.checkCharacters) > 0 {
		output = append(output, c.checkCharacters...)
	}

	return output
}

func elementsInArray(checkElements []int, char int, elements []int) int {
	i := 0
	offset := 0

	if char != -1 {
		checkElements = append(checkElements, char)
	}

	for _, element := range elements {

		if len(checkElements) <= offset {
			return i - len(checkElements)
		}

		if checkElements[offset] == element {
			offset += 1
		} else {
			offset = 0
		}

		i += 1
	}

	return -1
}
