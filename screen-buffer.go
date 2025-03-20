package blitra

import (
	"fmt"
	"io"
	"strconv"
)

const (
	escSetFGTrueColor = "\x1b[38;2;%d;%d;%dm"
	escSetBGTrueColor = "\x1b[48;2;%d;%d;%dm"
	escSetFG8Color    = "\x1b[3%dm"
	escSetBG8Color    = "\x1b[4%dm"
	escResetFGColor   = "\x1b[39m"
	escResetBGColor   = "\x1b[49m"
)

type ScreenBuffer struct {
	X               int
	Y               int
	Width           int
	Height          int
	Cells           []ScreenCell
	PrevCells       []ScreenCell
	TargetTTYStdout io.Writer
}

func NewScreenBuffer(x, y, width, height int, targetTTYStdout io.Writer) *ScreenBuffer {
	return &ScreenBuffer{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Cells:           make([]ScreenCell, width*height),
		PrevCells:       make([]ScreenCell, width*height),
		TargetTTYStdout: targetTTYStdout,
	}
}

func (sb *ScreenBuffer) MaybeResize(x, y, width, height int) {
	if sb.X == x && sb.Y == y && sb.Width == width && sb.Height == height {
		return
	}
	sb.X = x
	sb.Y = y
	sb.Width = width
	sb.Height = height
	sb.Cells = make([]ScreenCell, width*height)
	sb.PrevCells = make([]ScreenCell, width*height)
}

func (sb *ScreenBuffer) Set(c, r int, cell ScreenCell, merge bool) {
	if c < 0 || c >= sb.Width || r < 0 || r >= sb.Height {
		return
	}
	if merge {
		sb.Cells[r*sb.Width+c].Merge(&cell)
	} else {
		sb.Cells[r*sb.Width+c] = cell
	}
}

func (sb *ScreenBuffer) Get(x, y int) (*ScreenCell, bool) {
	if x < 0 || x >= sb.Width || y < 0 || y >= sb.Height {
		return &ScreenCell{}, false
	}
	cell := &sb.Cells[y*sb.Width+x]
	prevCell := &sb.PrevCells[y*sb.Width+x]
	return cell, !cell.IsEqual(prevCell)
}

func (sb *ScreenBuffer) DrawFrame() {
	width := sb.Width
	height := sb.Height
	x := sb.X
	y := sb.Y

	prevFgColor := ""
	prevBgColor := ""

	for r := 0; r < height; r += 1 {
		for c := 0; c < width; c += 1 {
			cell, isDirty := sb.Get(x+c, y+r)
			if !isDirty {
				continue
			}

			fmt.Fprintf(sb.TargetTTYStdout, escMoveCursor, y+r+1, x+c+1)

			// Set colors
			if cell.ForegroundColor != nil {
				fgColor := toForegroundColorEsc(*cell.ForegroundColor)
				if fgColor != prevFgColor {
					fmt.Fprint(sb.TargetTTYStdout, fgColor)
				}
				prevFgColor = fgColor
			}
			if cell.BackgroundColor != nil {
				bgColor := toBackgroundColorEsc(*cell.BackgroundColor)
				if bgColor != prevBgColor {
					fmt.Fprint(sb.TargetTTYStdout, bgColor)
				}
				prevBgColor = bgColor
			}

			// Set character
			fmt.Fprint(sb.TargetTTYStdout, string(VOr(cell.Character, ' ')))
		}
	}

	sb.PrevCells = make([]ScreenCell, width*height)
	for i, cell := range sb.Cells {
		sb.PrevCells[i].Merge(&cell)
	}
}

type ScreenCell struct {
	Character       *rune
	ForegroundColor *string
	BackgroundColor *string
	Bold            *bool
	Dim             *bool
	Italic          *bool
	Underline       *bool
	Blink           *bool
	FastBlink       *bool
	Hidden          *bool
	StrikeThrough   *bool
	DoubleUnderline *bool
}

func (sc *ScreenCell) IsEqual(other *ScreenCell) bool {
	if (sc.Character != nil && other.Character == nil) ||
		(sc.Character == nil && other.Character != nil) ||
		(sc.Character != nil && other.Character != nil && *sc.Character != *other.Character) {
		return false
	}
	if (sc.ForegroundColor != nil && other.ForegroundColor == nil) ||
		(sc.ForegroundColor == nil && other.ForegroundColor != nil) ||
		(sc.ForegroundColor != nil && other.ForegroundColor != nil && *sc.ForegroundColor != *other.ForegroundColor) {
		return false
	}
	if (sc.BackgroundColor != nil && other.BackgroundColor == nil) ||
		(sc.BackgroundColor == nil && other.BackgroundColor != nil) ||
		(sc.BackgroundColor != nil && other.BackgroundColor != nil && *sc.BackgroundColor != *other.BackgroundColor) {
		return false
	}
	if (sc.Bold != nil && other.Bold == nil) ||
		(sc.Bold == nil && other.Bold != nil) ||
		(sc.Bold != nil && other.Bold != nil && *sc.Bold != *other.Bold) {
		return false
	}
	if (sc.Dim != nil && other.Dim == nil) ||
		(sc.Dim == nil && other.Dim != nil) ||
		(sc.Dim != nil && other.Dim != nil && *sc.Dim != *other.Dim) {
		return false
	}
	if (sc.Italic != nil && other.Italic == nil) ||
		(sc.Italic == nil && other.Italic != nil) ||
		(sc.Italic != nil && other.Italic != nil && *sc.Italic != *other.Italic) {
		return false
	}
	if (sc.Underline != nil && other.Underline == nil) ||
		(sc.Underline == nil && other.Underline != nil) ||
		(sc.Underline != nil && other.Underline != nil && *sc.Underline != *other.Underline) {
		return false
	}
	if (sc.Blink != nil && other.Blink == nil) ||
		(sc.Blink == nil && other.Blink != nil) ||
		(sc.Blink != nil && other.Blink != nil && *sc.Blink != *other.Blink) {
		return false
	}
	if (sc.FastBlink != nil && other.FastBlink == nil) ||
		(sc.FastBlink == nil && other.FastBlink != nil) ||
		(sc.FastBlink != nil && other.FastBlink != nil && *sc.FastBlink != *other.FastBlink) {
		return false
	}
	if (sc.Hidden != nil && other.Hidden == nil) ||
		(sc.Hidden == nil && other.Hidden != nil) ||
		(sc.Hidden != nil && other.Hidden != nil && *sc.Hidden != *other.Hidden) {
		return false
	}
	if (sc.StrikeThrough != nil && other.StrikeThrough == nil) ||
		(sc.StrikeThrough == nil && other.StrikeThrough != nil) ||
		(sc.StrikeThrough != nil && other.StrikeThrough != nil && *sc.StrikeThrough != *other.StrikeThrough) {
		return false
	}
	if (sc.DoubleUnderline != nil && other.DoubleUnderline == nil) ||
		(sc.DoubleUnderline == nil && other.DoubleUnderline != nil) ||
		(sc.DoubleUnderline != nil && other.DoubleUnderline != nil && *sc.DoubleUnderline != *other.DoubleUnderline) {
		return false
	}
	return true
}

func (sc *ScreenCell) Merge(other *ScreenCell) {
	if other.Character != nil {
		character := *other.Character
		sc.Character = &character
	}
	if other.ForegroundColor != nil {
		fg := *other.ForegroundColor
		sc.ForegroundColor = &fg
	}
	if other.BackgroundColor != nil {
		bg := *other.BackgroundColor
		sc.BackgroundColor = &bg
	}
	if other.Bold != nil {
		bold := *other.Bold
		sc.Bold = &bold
	}
	if other.Dim != nil {
		dim := *other.Dim
		sc.Dim = &dim
	}
	if other.Italic != nil {
		italic := *other.Italic
		sc.Italic = &italic
	}
	if other.Underline != nil {
		underline := *other.Underline
		sc.Underline = &underline
	}
	if other.Blink != nil {
		blink := *other.Blink
		sc.Blink = &blink
	}
	if other.FastBlink != nil {
		fastBlink := *other.FastBlink
		sc.FastBlink = &fastBlink
	}
	if other.Hidden != nil {
		hidden := *other.Hidden
		sc.Hidden = &hidden
	}
	if other.StrikeThrough != nil {
		strikeThrough := *other.StrikeThrough
		sc.StrikeThrough = &strikeThrough
	}
	if other.DoubleUnderline != nil {
		doubleUnderline := *other.DoubleUnderline
		sc.DoubleUnderline = &doubleUnderline
	}
}

func toForegroundColorEsc(color string) string {
	if isRgbColor(color) {
		red, green, blue := toRgbColor(color)
		return fmt.Sprintf(escSetFGTrueColor, red, green, blue)
	}
	if is8Color(color) {
		return fmt.Sprintf(escSetFG8Color, to8Color(color))
	}
	return escResetFGColor
}

func toBackgroundColorEsc(color string) string {
	if isRgbColor(color) {
		red, green, blue := toRgbColor(color)
		return fmt.Sprintf(escSetBGTrueColor, red, green, blue)
	}
	if is8Color(color) {
		return fmt.Sprintf(escSetBG8Color, to8Color(color))
	}
	return escResetBGColor
}

func isRgbColor(color string) bool {
	c := []rune(color)
	len := len(c)
	return c[0] == '#' && len == 4 || len == 7
}

func toRgbColor(color string) (int, int, int) {
	c := []rune(color)
	len := len(c)
	if len == 4 {
		c = []rune{c[0], c[1], c[1], c[2], c[2], c[3], c[3]}
	}

	redStr := string(c[1:3])
	red, err := strconv.ParseInt(redStr, 16, 16)
	if err != nil {
		red = 0
	}

	greenStr := string(c[3:5])
	green, err := strconv.ParseInt(greenStr, 16, 16)
	if err != nil {
		green = 0
	}

	blueStr := string(c[5:7])
	blue, err := strconv.ParseInt(blueStr, 16, 16)
	if err != nil {
		blue = 0
	}

	return int(red), int(green), int(blue)
}

func is8Color(color string) bool {
	switch color {
	case "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white":
		return true
	default:
		return false
	}
}

func to8Color(color string) int {
	switch color {
	case "black":
		return 0
	case "red":
		return 1
	case "green":
		return 2
	case "yellow":
		return 3
	case "blue":
		return 4
	case "magenta":
		return 5
	case "cyan":
		return 6
	case "white":
		return 7
	default:
		panic("invalid 8-color")
	}
}
