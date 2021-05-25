package lzss

import (
	"fmt"
	"io"

	"github.com/xinmyname/bitstream-go"
)

type compressor struct {
	bs        bitstream.Writer
	windowLen int
}

func (c *compressor) Write(p []byte) (n int, err error) {

	compressed := compress(p, c.windowLen)

	for _, token := range compressed {
		if token > 255 {
			length := token & 0xffff
			offset := token >> 16
			c.bs.WriteBit(1)
			c.bs.WriteUint16BE(uint16(offset))
			c.bs.WriteUint16BE(uint16(length))
		} else {
			c.bs.WriteBit(0)
			c.bs.WriteUint8(uint8(token))
		}
	}

	return len(p), err
}

func (c *compressor) Close() error {
	return c.bs.Flush()
}

func NewWriter(w io.Writer) io.WriteCloser {
	return NewWriterWindow(w, 256)
}

func NewWriterWindow(w io.Writer, windowLen int) io.WriteCloser {

	return &compressor{
		bs:        *bitstream.NewWriter(w),
		windowLen: windowLen,
	}
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

func compress(textBytes []byte, maxSlidingWindowSize int) []int {
	searchBuffer := make([]int, 0, maxSlidingWindowSize)
	checkCharacters := make([]int, 0, maxSlidingWindowSize)
	output := make([]int, 0, maxSlidingWindowSize)
	i := 0
	movedPast := 0

	for _, char := range textBytes {

		index := elementsInArray(checkCharacters, int(char), searchBuffer)

		if index == -1 || i == len(textBytes)-1 {
			if i == len(textBytes)-1 && index != -1 {
				checkCharacters = append(checkCharacters, int(char))
			}

			if len(checkCharacters) > 1 {
				index = elementsInArray(checkCharacters, -1, searchBuffer)
				offset := i - index - len(checkCharacters) - movedPast
				length := len(checkCharacters)
				token := ((offset & 0x7fff) << 16) | (length & 0xffff)

				pyToken := fmt.Sprintf("<%d,%d>", offset, length)
				lenToken := len(pyToken)

				if lenToken > length {
					output = append(output, checkCharacters...)
				} else {
					output = append(output, token)
				}

				searchBuffer = append(searchBuffer, checkCharacters...)

			} else {
				output = append(output, checkCharacters...)
				searchBuffer = append(searchBuffer, checkCharacters...)
			}

			checkCharacters = checkCharacters[:0]
		}

		checkCharacters = append(checkCharacters, int(char))

		if len(searchBuffer) > maxSlidingWindowSize {
			diff := len(searchBuffer) - maxSlidingWindowSize
			movedPast += diff
			searchBuffer = searchBuffer[diff:]
		}

		i += 1
	}

	return output
}
