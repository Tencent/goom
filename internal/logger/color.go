package logger

import "fmt"

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
