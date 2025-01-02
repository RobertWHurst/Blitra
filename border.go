package blitra

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

func NewBorder(topLeft, topRight, bottomLeft, bottomRight, left, right, top, bottom string) *Border {
	return &Border{
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomLeft:  bottomLeft,
		BottomRight: bottomRight,
		Left:        left,
		Right:       right,
		Top:         top,
		Bottom:      bottom,
	}
}

func DoubleBorder() *Border {
	return NewBorder("╔", "╗", "╚", "╝", "║", "║", "═", "═")
}

func RoundBorder() *Border {
	return NewBorder("╭", "╮", "╰", "╯", "│", "│", "─", "─")
}

func BoldBorder() *Border {
	return NewBorder("┏", "┓", "┗", "┛", "┃", "┃", "━", "━")
}

func LightBorder() *Border {
	return NewBorder("┌", "┐", "└", "┘", "│", "│", "─", "─")
}
