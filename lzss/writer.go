package lzss

import (
	"bufio"
	"fmt"
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

	compressed := compress(p, 256)

	for _, c := range compressed {
		if c > 255 {
			length := c & 0xffff
			offset := c >> 16

			fmt.Printf("<%d,%d>", offset, length)
		} else {
			fmt.Printf("%c", rune(c))
		}
	}

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

func elementsInArrayPlusChar(checkElements []int, char int, elements []int) int {
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

func elementsInArray(checkElements []int, elements []int) int {
	return elementsInArrayPlusChar(checkElements, -1, elements)
}

func compress(textBytes []byte, maxSlidingWindowSize int) []int {
	searchBuffer := make([]int, 0, maxSlidingWindowSize)
	checkCharacters := make([]int, 0, maxSlidingWindowSize)
	output := make([]int, 0, maxSlidingWindowSize)
	i := 0
	movedPast := 0

	for _, char := range textBytes {

		if elementsInArrayPlusChar(checkCharacters, int(char), searchBuffer) == -1 || i == len(textBytes)-1 {
			if i == len(textBytes)-1 && elementsInArrayPlusChar(checkCharacters, int(char), searchBuffer) != -1 {
				checkCharacters = append(checkCharacters, int(char))
			}

			if len(checkCharacters) > 1 {
				index := elementsInArray(checkCharacters, searchBuffer)
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
