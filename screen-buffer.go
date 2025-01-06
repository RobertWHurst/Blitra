package blitra

type ScreenBuffer struct {
	X         int
	Y         int
	Width     int
	Height    int
	Cells     []ScreenCell
	PrevCells []ScreenCell
}

func NewScreenBuffer(x, y, width, height int) *ScreenBuffer {
	return &ScreenBuffer{
		X:         x,
		Y:         y,
		Width:     width,
		Height:    height,
		Cells:     make([]ScreenCell, width*height),
		PrevCells: make([]ScreenCell, width*height),
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

func (sb *ScreenBuffer) Set(x, y int, cell ScreenCell) {
	if x < 0 || x >= sb.Width || y < 0 || y >= sb.Height {
		return
	}
	sb.Cells[y*sb.Width+x] = cell
}

func (sb *ScreenBuffer) Get(x, y int) (*ScreenCell, bool) {
	if x < 0 || x >= sb.Width || y < 0 || y >= sb.Height {
		return &ScreenCell{}, false
	}
	cell := &sb.Cells[y*sb.Width+x]
	prevCell := &sb.PrevCells[y*sb.Width+x]
	return cell, !cell.IsEqual(prevCell)
}

func (sb *ScreenBuffer) MarkFrame() {
	copy(sb.PrevCells, sb.Cells)
}

type ScreenCell struct {
	Character       rune
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
	return sc.Character == other.Character &&
		compCellP(sc.ForegroundColor, other.ForegroundColor) &&
		compCellP(sc.BackgroundColor, other.BackgroundColor) &&
		compCellP(sc.Bold, other.Bold) &&
		compCellP(sc.Dim, other.Dim) &&
		compCellP(sc.Italic, other.Italic) &&
		compCellP(sc.Underline, other.Underline) &&
		compCellP(sc.Blink, other.Blink) &&
		compCellP(sc.FastBlink, other.FastBlink) &&
		compCellP(sc.Hidden, other.Hidden) &&
		compCellP(sc.StrikeThrough, other.StrikeThrough) &&
		compCellP(sc.DoubleUnderline, other.DoubleUnderline)
}

func compCellP[T string | bool](a, b *T) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return *a == *b
}
