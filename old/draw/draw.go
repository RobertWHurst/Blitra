package draw

import (
	"strings"

	"github.com/robertwhurst/blitra/cursor"
	"github.com/robertwhurst/blitra/stdio"
)

func AtCursor(str string) {
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		stdio.Print(line)
		cursor.Down(1)
		lineRuneLen := len([]rune(line))
		cursor.Left(lineRuneLen)
	}
}

func AtPosition(row int, col int, str string) {
	cursor.Position(row, col)
	stdio.Print(str)
}
