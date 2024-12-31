package blitra

type ScreenBuffer struct {
	X       int
	Y       int
	Width   int
	Height  int
	Content []ScreenCell
}

func NewScreenBuffer(x, y, width, height int) *ScreenBuffer {
	return &ScreenBuffer{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Content: make([]ScreenCell, width*height),
	}
}

func (sb *ScreenBuffer) MaybeResize(x, y, width, height int) {
	if x != sb.X || y != sb.Y || width != sb.Width || height != sb.Height {
		sb.X = x
		sb.Y = y
		sb.Width = width
		sb.Height = height
		oldContent := sb.Content
		sb.Content = make([]ScreenCell, width*height)
		for r := 0; r < height; r += 1 {
			for c := 0; c < width; c += 1 {
				if r < sb.Height && c < sb.Width {
					sb.Content[r*width+c] = oldContent[r*sb.Width+c]
				}
			}
		}
	}
}

func (sb *ScreenBuffer) Set(x, y int, cell ScreenCell) {
	if x < 0 || x >= sb.Width || y < 0 || y >= sb.Height {
		return
	}
	sb.Content[y*sb.Width+x] = cell
}

func (sb *ScreenBuffer) Get(x, y int) ScreenCell {
	if x < 0 || x >= sb.Width || y < 0 || y >= sb.Height {
		return ScreenCell{}
	}
	return sb.Content[y*sb.Width+x]
}

type ScreenCell struct {
	Character       rune
	ForegroundColor *string
	BackgroundColor *string
	// TODO: The following properties are not yet implemented.
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
