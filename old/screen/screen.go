// This package contains functions for manipulating the state of the
// terminal screen. Most notably, it contains functions for clearing the
// screen, clearing lines, and switching between the main and alternate
// screen buffers.
package screen

import (
	"fmt"

	"github.com/robertwhurst/blitra/stdio"
)

// Clear the visible area of the terminal screen and scrollback history.
func Clear() {
	fmt.Print("\x1B[2J")
}

// Clear the visible area of the terminal screen at and below the cursor.
func ClearDown() {
	fmt.Print("\x1B[J")
}

// Clear the visible area of the terminal screen at and above the cursor.
func ClearUp() {
	fmt.Print("\x1B[1J")
}

// Clear the current line of the terminal screen.
func ClearLine() {
	fmt.Print("\x1B[2K")
}

// Clear the current line of the terminal screen to the right and including
// the current cursor column.
func ClearLineToRight() {
	fmt.Print("\x1B[K")
}

// Clear the current line of the terminal screen to the left and including
// the current cursor column.
func ClearLineToLeft() {
	fmt.Print("\x1B[1K")
}

// Clear the scrollback history of the terminal screen, leaving the visible
// area intact.
func ClearScrollback() {
	fmt.Print("\x1B[3J")
}

// Switch to the alternate screen buffer.
func AlternateBuffer() {
	stdio.Print("\x1B[?1049h")
}

// Switch to the main screen buffer.
func MainBuffer() {
	stdio.Print("\x1B[?1049l")
}
