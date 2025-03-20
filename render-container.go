package blitra

func renderContainer(el *Element, screenBuffer *ScreenBuffer) error {
	x := el.Position.X + el.LeftMargin()
	y := el.Position.Y + el.TopMargin()
	w := max(el.Size.Width-el.HorizontalMargin(), 0)
	h := max(el.Size.Height-el.VerticalMargin(), 0)
	fg := el.Style.TextColor
	bg := el.Style.BackgroundColor

	// foreground and background color
	for r := y; r < y+h; r += 1 {
		for c := x; c < x+w; c += 1 {
			screenBuffer.Set(c, r, ScreenCell{
				ForegroundColor: fg,
				BackgroundColor: bg,
				// TODO: Implement these
				// Bold:            nil,
				// Dim:             nil,
				// Italic:          nil,
				// Underline:       nil,
				// Blink:           nil,
				// FastBlink:       nil,
				// Hidden:          nil,
				// StrikeThrough:   nil,
				// DoubleUnderline: nil,
			}, el.Parent != nil)
		}
	}

	// border

	return nil
}
