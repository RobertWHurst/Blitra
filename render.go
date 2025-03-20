package blitra

import (
	"fmt"
	"strings"
)

const (
	escEnableMouse  = "\x1b[?1000h\x1b[?1003h\x1b[?1006h"
	escDisableMouse = "\x1b[?1006l\x1b[?1003l\x1b[?1000l"

	escEnableFocusTracking  = "\x1b[?1004h"
	escDisableFocusTracking = "\x1b[?1004l"

	escHideCursor = "\x1b[?25l"
	escShowCursor = "\x1b[?25h"

	escSecondaryScreen = "\x1b[?1049h"
	escPrimaryScreen   = "\x1b[?1049l"

	escMoveCursor = "\x1b[%d;%dH"
	escErase      = "\x1b[%dX"
)

var DebugDraw = false

func PrepareScreen(view *ViewHandle) {
	fmt.Fprint(view.stdioManager.targetTTYStdout, escHideCursor)
	fmt.Fprint(view.stdioManager.targetTTYStdout, escEnableMouse)
	fmt.Fprint(view.stdioManager.targetTTYStdout, escEnableFocusTracking)
	if view.opts.TargetBuffer == SecondaryBuffer {
		fmt.Fprint(view.stdioManager.targetTTYStdout, escSecondaryScreen)
	}
}

func RestoreScreen(view *ViewHandle) {
	if view.opts.TargetBuffer == SecondaryBuffer {
		fmt.Fprint(view.stdioManager.targetTTYStdout, escPrimaryScreen)
	}
	fmt.Fprint(view.stdioManager.targetTTYStdout, escDisableFocusTracking)
	fmt.Fprint(view.stdioManager.targetTTYStdout, escDisableMouse)
	fmt.Fprint(view.stdioManager.targetTTYStdout, escShowCursor)
}

func Render(view *ViewHandle, rootElement *Element) error {
	if err := VisitElementsDown(rootElement, view.screenBuffer, renderElementVisitor); err != nil {
		return fmt.Errorf("Failed to render element tree: %w", err)
	}
	view.screenBuffer.DrawFrame()

	return nil
}

func renderElementVisitor(el *Element, screenBuffer *ScreenBuffer) error {
	var err error
	switch el.Kind {
	case TextElementKind:
		err = renderText(el, screenBuffer)
	case ContainerElementKind:
		err = renderContainer(el, screenBuffer)
	}

	if DebugDraw {
		screenBuffer.DrawFrame()
	}

	return err
}

// func xrenderBorders(view *ViewHandle, element *Element) {
// 	sb := view.screenBuffer

// 	x := element.Position.X + V(element.Style.LeftMargin)
// 	y := element.Position.Y + V(element.Style.TopMargin)
// 	width := element.Size.Width - V(element.Style.LeftMargin) - V(element.Style.RightMargin)
// 	height := element.Size.Height - V(element.Style.TopMargin) - V(element.Style.BottomMargin)

// 	if element.Style.LeftBorder != nil {
// 		bCells, bWidth, bHeight := strToCells(element.Style.LeftBorder.left, element.Style, -1, -1)
// 		for r := 0; r < height; r += 1 {
// 			for c := 0; c < bWidth; c += 1 {
// 				rr := c % bHeight
// 				sb.Set(x+c, y+r, bCells[rr][c])
// 			}
// 		}
// 	}
// 	if element.Style.RightBorder != nil {
// 		bCells, bWidth, bHeight := strToCells(element.Style.RightBorder.right, element.Style, -1, -1)
// 		for r := 0; r < height; r += 1 {
// 			for c := 0; c < bWidth; c += 1 {
// 				rr := c % bHeight
// 				sb.Set(x+width-bWidth+c, y+r, bCells[rr][c])
// 			}
// 		}
// 	}
// 	if element.Style.TopBorder != nil {
// 		lBCells, lBWidth, _ := strToCells(element.Style.TopBorder.topLeft, element.Style, -1, -1)
// 		rBCells, rBWidth, _ := strToCells(element.Style.TopBorder.topRight, element.Style, -1, -1)
// 		bCells, bWidth, bHeight := strToCells(element.Style.TopBorder.top, element.Style, -1, -1)
// 		for r := 0; r < bHeight; r += 1 {
// 			for c := 0; c < width; c += 1 {
// 				if c < lBWidth {
// 					cc := c % lBWidth
// 					sb.Set(x+c, y+r, lBCells[r][cc])
// 				} else if c >= width-rBWidth {
// 					cc := c % rBWidth
// 					sb.Set(x+c, y+r, rBCells[r][cc])
// 				} else {
// 					cc := c % bWidth
// 					sb.Set(x+c, y+r, bCells[r][cc])
// 				}
// 			}
// 		}
// 	}
// 	if element.Style.BottomBorder != nil {
// 		lBCells, lBWidth, _ := strToCells(element.Style.BottomBorder.bottomLeft, element.Style, -1, -1)
// 		rBCells, rBWidth, _ := strToCells(element.Style.BottomBorder.bottomRight, element.Style, -1, -1)
// 		bCells, bWidth, bHeight := strToCells(element.Style.BottomBorder.bottom, element.Style, -1, -1)
// 		for r := 0; r < bHeight; r += 1 {
// 			for c := 0; c < width; c += 1 {
// 				if c < lBWidth {
// 					cc := c % lBWidth
// 					sb.Set(x+c, y+height-bHeight+r, lBCells[r][cc])
// 				} else if c >= width-rBWidth {
// 					cc := c % rBWidth
// 					sb.Set(x+c, y+height-bHeight+r, rBCells[r][cc])
// 				} else {
// 					cc := c % bWidth
// 					sb.Set(x+c, y+height-bHeight+r, bCells[r][cc])
// 				}
// 			}
// 		}
// 	}
// }

// func xrenderText(view *ViewHandle, element *Element) {
// 	sb := view.screenBuffer
// 	x := element.Position.X + V(element.Style.LeftMargin)
// 	y := element.Position.Y + V(element.Style.TopMargin)
// 	w := element.Size.Width - V(element.Style.LeftMargin) - V(element.Style.RightMargin)
// 	h := element.Size.Height - V(element.Style.TopMargin) - V(element.Style.BottomMargin)

// 	text, width, height := strToCells(element.Text, element.Style, w, h)
// 	for r := 0; r < height; r += 1 {
// 		for c := 0; c < width; c += 1 {
// 			sb.Set(x+c, y+r, text[r][c])
// 		}
// 	}
// }

func strToCells(s string, style Style, width, height int) ([][]ScreenCell, int, int) {
	lines := strings.Split(s, "\n")
	sHeight := len(lines)
	sWidth := 0
	for _, line := range lines {
		lineLen := len([]rune(line))
		if lineLen > sWidth {
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
				Character:       &char,
				ForegroundColor: style.TextColor,
				BackgroundColor: style.BackgroundColor,
			}
		}
	}

	return cells, width, height
}
