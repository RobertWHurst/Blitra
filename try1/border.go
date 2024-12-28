package tui

type Border struct {
	TopLeft     string
	Top         string
	TopRight    string
	BottomLeft  string
	Bottom      string
	BottomRight string
	Left        string
	Right       string
}

func NewBorder(topLeft, top, topRight, bottomLeft, bottom, bottomRight, left, right string) Border {
	return Border{
		TopLeft:     topLeft,
		Top:         top,
		TopRight:    topRight,
		BottomLeft:  bottomLeft,
		Bottom:      bottom,
		BottomRight: bottomRight,
		Left:        left,
		Right:       right,
	}
}

func DoubleBorder() Border {
	return NewBorder("╔", "═", "╗", "╚", "═", "╝", "║", "║")
}

func RoundBorder() Border {
	return NewBorder("╭", "─", "╮", "╰", "─", "╯", "│", "│")
}

func BoldBorder() Border {
	return NewBorder("┏", "━", "┓", "┗", "━", "┛", "┃", "┃")
}

func LightBorder() Border {
	return NewBorder("┌", "─", "┐", "└", "─", "┘", "│", "│")
}
