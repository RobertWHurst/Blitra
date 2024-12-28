package main

import (
	"time"

	"github.com/robertwhurst/blitra/box"
	"github.com/robertwhurst/blitra/cursor"
	"github.com/robertwhurst/blitra/draw"
	"github.com/robertwhurst/blitra/format"
	"github.com/robertwhurst/blitra/screen"
)

// # Colors - 3/4 bit (16 colors)
// # Foreground colors
// FG_BLACK = "\033[30m"
// FG_RED = "\033[31m"
// FG_GREEN = "\033[32m"
// FG_YELLOW = "\033[33m"
// FG_BLUE = "\033[34m"
// FG_MAGENTA = "\033[35m"
// FG_CYAN = "\033[36m"
// FG_WHITE = "\033[37m"
// # Bright foreground colors
// FG_BRIGHT_BLACK = "\033[90m"
// FG_BRIGHT_RED = "\033[91m"
// FG_BRIGHT_GREEN = "\033[92m"
// FG_BRIGHT_YELLOW = "\033[93m"
// FG_BRIGHT_BLUE = "\033[94m"
// FG_BRIGHT_MAGENTA = "\033[95m"
// FG_BRIGHT_CYAN = "\033[96m"
// FG_BRIGHT_WHITE = "\033[97m"

// # Background colors
// BG_BLACK = "\033[40m"
// BG_RED = "\033[41m"
// BG_GREEN = "\033[42m"
// BG_YELLOW = "\033[43m"
// BG_BLUE = "\033[44m"
// BG_MAGENTA = "\033[45m"
// BG_CYAN = "\033[46m"
// BG_WHITE = "\033[47m"
// # Bright background colors
// BG_BRIGHT_BLACK = "\033[100m"
// BG_BRIGHT_RED = "\033[101m"
// BG_BRIGHT_GREEN = "\033[102m"
// BG_BRIGHT_YELLOW = "\033[103m"
// BG_BRIGHT_BLUE = "\033[104m"
// BG_BRIGHT_MAGENTA = "\033[105m"
// BG_BRIGHT_CYAN = "\033[106m"
// BG_BRIGHT_WHITE = "\033[107m"

// # 256 colors and RGB
// FG_256_COLOR = "\033[38;5;%dm"         # Set foreground color (256 colors)
// BG_256_COLOR = "\033[48;5;%dm"         # Set background color (256 colors)
// FG_RGB_COLOR = "\033[38;2;%d;%d;%dm"   # Set foreground color (RGB)
// BG_RGB_COLOR = "\033[48;2;%d;%d;%dm"   # Set background color (RGB)

func main() {
	// row, col := cursor.MustGetPosition()

	screen.AlternateBuffer()
	screen.Clear()
	cursor.Hide()
	cursor.Position(10, 10)
	// box.Draw(10, 10, 15, 3)

	cursor.Position(11, 11)
	format.Bold()
	draw.AtCursor(box.Wrap("Hello, World!"))
	cursor.Down(10)
	draw.AtCursor(box.Wrap("Hello, World, again!"))

	time.Sleep(10 * time.Second)

	screen.MainBuffer()
	cursor.Show()
	// cursor.Position(row, col)
}
