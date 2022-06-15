package logger

import (
	"bufio"
	"fmt"
	"strings"
)

// Foreground colors.
const (
	None Color = iota + 29
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	if c == None {
		return s
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

// AddAll adds the coloring to the given string of all line.
func (c Color) AddAll(s string) string {
	if c == None {
		return s
	}
	lines := make([]string, 0, 10)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, c.Add(scanner.Text()))
	}
	return strings.Join(lines, "\n")
}
