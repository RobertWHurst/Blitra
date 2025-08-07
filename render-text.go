package blitra

func renderText(el *Element, screenBuffer *ScreenBuffer) error {
	contentX := el.Position.X + el.LeftEdge()
	contentY := el.Position.Y + el.TopEdge()

	contentWidth := el.Size.Width - el.HorizontalEdge()
	contentHeight := el.Size.Height - el.VerticalEdge()

	if contentWidth <= 0 || contentHeight <= 0 {
		return nil
	}

	textCells, textWidth, textHeight := strToCells(
		el.Text,
		*calcTextStyle(el),
		contentWidth,
		contentHeight,
	)

	for r := 0; r < textHeight; r++ {
		for c := 0; c < textWidth; c++ {
			if r < len(textCells) && c < len(textCells[r]) {
				screenBuffer.Set(contentX+c, contentY+r, textCells[r][c], el.Parent != nil)
			}
		}
	}

	return nil
}

func calcTextStyle(el *Element) *Style {
	textColor := el.Style.TextColor

	tEl := el
	for textColor == nil && tEl != nil {
		tEl = tEl.Parent
		textColor = tEl.Style.TextColor
	}

	return &Style{
		TextColor: textColor,
	}
}
