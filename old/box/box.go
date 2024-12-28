package box

import (
	"strings"

	"github.com/robertwhurst/blitra/cursor"
	"github.com/robertwhurst/blitra/stdio"
)

var (
	UpperLeftCorner  = "┌"
	UpperRightCorner = "┐"
	LowerLeftCorner  = "└"
	LowerRightCorner = "┘"
	HorizontalLine   = "─"
	VerticalLine     = "│"
	LeftTee          = "├"
	RightTee         = "┤"
	Cross            = "┼"
	TopTee           = "┬"
	BottomTee        = "┴"
)

func Draw(startCol, startRow, colSpan, rowSpan int) {
	endCol := startCol + (colSpan - 1)
	endRow := startRow + (rowSpan - 1)

	stdio.Print(UpperLeftCorner)
	for i := startCol + 1; i < endCol; i += 1 {
		stdio.Print(HorizontalLine)
	}
	stdio.Print(UpperRightCorner)

	for i := startRow + 1; i < endRow; i += 1 {
		cursor.Down(1)
		cursor.Column(startCol)
		stdio.Print(VerticalLine)
		cursor.Column(endCol)
		stdio.Print(VerticalLine)
	}

	cursor.Down(1)
	cursor.Column(startCol)
	stdio.Print(LowerLeftCorner)
	for i := startCol + 1; i < endCol; i += 1 {
		stdio.Print(HorizontalLine)
	}
	stdio.Print(LowerRightCorner)
}

func Wrap(str string) string {
	lines := strings.Split(str, "\n")
	lineLen := 0
	for _, line := range lines {
		if len(line) > lineLen {
			lineLen = len(line)
		}
	}

	wrapped := UpperLeftCorner
	for i := 0; i < lineLen; i += 1 {
		wrapped += HorizontalLine
	}
	wrapped += UpperRightCorner + "\n"

	for _, line := range lines {
		wrapped += VerticalLine + padLine(line, lineLen) + VerticalLine + "\n"
	}

	wrapped += LowerLeftCorner
	for i := 0; i < lineLen; i += 1 {
		wrapped += HorizontalLine
	}
	wrapped += LowerRightCorner

	return wrapped
}

func padLine(line string, length int) string {
	return line + strings.Repeat(" ", length-len(line))
}
