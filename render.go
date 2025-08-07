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
