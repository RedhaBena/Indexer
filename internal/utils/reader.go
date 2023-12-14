package utils

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type CustomReader struct {
	before      [19]byte // max int is 2147483647 --> 9 letter + 11 for numberint() = 9+11 = 20 -1 = 19
	innerReader *bufio.Reader
}

func NewCustomReader(reader io.Reader) *CustomReader {
	return &CustomReader{
		innerReader: bufio.NewReader(reader),
	}
}

func (cr *CustomReader) Read(p []byte) (n int, err error) {
	// Read from the inner reader
	n, err = cr.innerReader.Read(p)
	if err != nil {
		return n, err
	}

	before := cr.before

	// Store the last 19 bytes in the 'before' field
	copy(cr.before[:], p[max(0, n-19):n])

	// Peek 19 bytes from the underlying reader
	peekBytes, _ := cr.innerReader.Peek(19)

	// Combine the last 19 bytes from the previous read, the current read, and the peeked bytes
	text := append(append(before[:], p[:n]...), peekBytes...)

	// Modify the data using removeNumberInt function
	modifiedContent := removeNumberInt(string(text))

	// Copy the modified content back to the original buffer
	copy(p[:n], []byte(modifiedContent)[len(before):len(modifiedContent)-len(peekBytes)])

	return n, err
}

func removeNumberInt(input string) string {
	pattern := `NumberInt\((\d+)\)`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Replace only the numeric part with spaces, keeping the rest
		numericPart := re.ReplaceAllString(match, "$1")
		return strings.Repeat(" ", 10) + numericPart + strings.Repeat(" ", 1)
	})
}

// Helper function to return the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
