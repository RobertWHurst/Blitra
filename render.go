package blitra

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	escHideCursor      = "\033[?25l"
	escShowCursor      = "\033[?25h"
	escSecondaryScreen = "\033[?1049h"
	escPrimaryScreen   = "\033[?1049l"
	escMoveCursor      = "\033[%d;%dH"
	escErase           = "\033[%dX"
	escSetFGTrueColor  = "\033[38;2;%d;%d;%dm"
	escSetBGTrueColor  = "\033[48;2;%d;%d;%dm"
	escResetFGColor    = "\033[39m"
	escResetBGColor    = "\033[49m"
)

func Bind(view *ViewHandle) {
	fmt.Fprint(view.tty, escHideCursor)
	if view.opts.TargetBuffer == SecondaryBuffer {
		fmt.Fprint(view.tty, escSecondaryScreen)
	}
}

func Unbind(view *ViewHandle) {
	if view.opts.TargetBuffer == SecondaryBuffer {
		fmt.Fprint(view.tty, escPrimaryScreen)
	}
	fmt.Fprint(view.tty, escShowCursor)
}

func Render(view *ViewHandle, element *Element) {
	element.VisitContainerElementsDown(func(el *Element) {
		renderElement(view, el)
	})
	renderView(view)
}

func renderElement(view *ViewHandle, element *Element) {
	renderText(view, element)
	renderBorders(view, element)
}

func renderBorders(view *ViewHandle, element *Element) {
	sb := view.screenBuffer
	x := element.Position.X
	y := element.Position.Y
	width := element.Size.Width
	height := element.Size.Height

	if element.Style.LeftBorder != nil {
		bCells, bWidth, bHeight := strToCells(element.Style.LeftBorder.Left, element.Style, -1, -1)
		for r := 0; r < height; r += 1 {
			for c := 0; c < bWidth; c += 1 {
				rr := c % bHeight
				sb.Set(x+c, y+r, bCells[rr][c])
			}
		}
	}
	if element.Style.RightBorder != nil {
		bCells, bWidth, bHeight := strToCells(element.Style.RightBorder.Right, element.Style, -1, -1)
		for r := 0; r < height; r += 1 {
			for c := 0; c < bWidth; c += 1 {
				rr := c % bHeight
				sb.Set(x+width-bWidth+c, y+r, bCells[rr][c])
			}
		}
	}
	if element.Style.TopBorder != nil {
		lBCells, lBWidth, _ := strToCells(element.Style.TopBorder.TopLeft, element.Style, -1, -1)
		rBCells, rBWidth, _ := strToCells(element.Style.TopBorder.TopRight, element.Style, -1, -1)
		bCells, bWidth, bHeight := strToCells(element.Style.TopBorder.Top, element.Style, -1, -1)
		for r := 0; r < bHeight; r += 1 {
			for c := 0; c < width; c += 1 {
				if c < lBWidth {
					cc := c % lBWidth
					sb.Set(x+c, y+r, lBCells[r][cc])
				} else if c >= width-rBWidth {
					cc := c % rBWidth
					sb.Set(x+c, y+r, rBCells[r][cc])
				} else {
					cc := c % bWidth
					sb.Set(x+c, y+r, bCells[r][cc])
				}
			}
		}
	}
	if element.Style.BottomBorder != nil {
		lBCells, lBWidth, _ := strToCells(element.Style.BottomBorder.BottomLeft, element.Style, -1, -1)
		rBCells, rBWidth, _ := strToCells(element.Style.BottomBorder.BottomRight, element.Style, -1, -1)
		bCells, bWidth, bHeight := strToCells(element.Style.BottomBorder.Bottom, element.Style, -1, -1)
		for r := 0; r < bHeight; r += 1 {
			for c := 0; c < width; c += 1 {
				if c < lBWidth {
					cc := c % lBWidth
					sb.Set(x+c, y+height-bHeight+r, lBCells[r][cc])
				} else if c >= width-rBWidth {
					cc := c % rBWidth
					sb.Set(x+c, y+height-bHeight+r, rBCells[r][cc])
				} else {
					cc := c % bWidth
					sb.Set(x+c, y+height-bHeight+r, bCells[r][cc])
				}
			}
		}
	}
}

func renderText(view *ViewHandle, element *Element) {
	sb := view.screenBuffer
	x := element.Position.X
	y := element.Position.Y
	w := element.Size.Width
	h := element.Size.Height

	text, width, height := strToCells(element.Text, element.Style, w, h)
	for r := 0; r < height; r += 1 {
		for c := 0; c < width; c += 1 {
			sb.Set(x+c, y+r, text[r][c])
		}
	}
}

func renderView(view *ViewHandle) {
	sb := view.screenBuffer
	x := view.x
	y := view.y
	width := view.width
	height := view.height

	prevFgColor := ""
	prevBgColor := ""

	for r := 0; r < height; r += 1 {
		for c := 0; c < width; c += 1 {
			cell, isDirty := sb.Get(x+c, y+r)
			if !isDirty {
				continue
			}

			fmt.Fprintf(view.tty, escMoveCursor, y+r+1, x+c+1)

			// debug marker
			// fmt.Fprint(view.tty, "X")
			// fmt.Fprintf(view.tty, escMoveCursor, y+r+1, x+c+1)

			// Set colors
			if cell.ForegroundColor != nil {
				fgColor := toForegroundColorEsc(*cell.ForegroundColor)
				if fgColor != prevFgColor {
					fmt.Fprint(view.tty, fgColor)
				}
				prevFgColor = fgColor
			}
			if cell.BackgroundColor != nil {
				bgColor := toBackgroundColorEsc(*cell.BackgroundColor)
				if bgColor != prevBgColor {
					fmt.Fprint(view.tty, bgColor)
				}
				prevBgColor = bgColor
			}

			// Set character
			fmt.Fprint(view.tty, string(cell.Character))
		}
	}

	sb.MarkFrame()
}

func strToCells(s string, style Style, width, height int) ([][]ScreenCell, int, int) {
	lines := strings.Split(s, "\n")
	sHeight := len(lines)
	sWidth := 0
	for _, line := range lines {
		lineLen := len([]rune(line))
		if lineLen > width {
			sWidth = lineLen
		}
	}

	if width == -1 {
		width = sWidth
	}
	if height == -1 {
		height = sHeight
	}

	cells := make([][]ScreenCell, height)
	for i := range cells {
		cells[i] = make([]ScreenCell, width)
	}

	for r := 0; r < height; r += 1 {
		for c := 0; c < width; c += 1 {
			char := ' '
			if r < len(lines) && c < len([]rune(lines[r])) {
				char = []rune(lines[r])[c]
			}
			cells[r][c] = ScreenCell{
				Character:       char,
				ForegroundColor: style.TextColor,
				BackgroundColor: style.BackgroundColor,
			}
		}
	}

	return cells, width, height
}

func toForegroundColorEsc(color string) string {
	c := []rune(color)
	len := len(c)
	if isRgbColor(c, len) {
		red, green, blue := toRgbColor(c, len)
		return fmt.Sprintf(escSetFGTrueColor, red, green, blue)
	}
	return escResetFGColor
}

func toBackgroundColorEsc(color string) string {
	c := []rune(color)
	len := len(c)
	if isRgbColor(c, len) {
		red, green, blue := toRgbColor(c, len)
		return fmt.Sprintf(escSetBGTrueColor, red, green, blue)
	}
	return escResetBGColor
}

func isRgbColor(c []rune, len int) bool {
	return c[0] == '#' && len == 4 || len == 7
}

func toRgbColor(c []rune, len int) (int, int, int) {
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

// func to8Color(color string) int {

// }
