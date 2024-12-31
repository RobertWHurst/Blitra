package blitra

import (
	"os"

	"golang.org/x/term"
)

// Determines if there is a tty. Checks stdout by default, but can check a
// a file if provided.
func IsTTY(tty *os.File) bool {
	return term.IsTerminal(int(tty.Fd()))
}

// Gets the size of the terminal. Checks stdout by default, but can check a
// a file if provided.
func MustGetTerminalSize(tty *os.File) (int, int) {
	width, height, err := term.GetSize(int(tty.Fd()))
	if err != nil {
		panic("Failed to get terminal size: " + err.Error())
	}
	return width, height
}
