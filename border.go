package blitra

import (
	"fmt"
	"math"
)

type Border struct {
	topLeft     string
	top         string
	topRight    string
	bottomLeft  string
	bottom      string
	bottomRight string
	left        string
	right       string

	topLeftCellSize     *Size
	topCellSize         *Size
	topRightCellSize    *Size
	leftCellSize        *Size
	rightCellSize       *Size
	bottomLeftCellSize  *Size
	bottomCellSize      *Size
	bottomRightCellSize *Size
}

func NewBorder(topLeft, topRight, bottomLeft, bottomRight, left, right, top, bottom string) *Border {
	return &Border{
		topLeft:     topLeft,
		topRight:    topRight,
		bottomLeft:  bottomLeft,
		bottomRight: bottomRight,
		left:        left,
		right:       right,
		top:         top,
		bottom:      bottom,
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

func (b *Border) TopLeftCellSize() Size {
	if b.topLeftCellSize == nil {
		size := getStrSize(b.topLeft)
		b.topLeftCellSize = &size
	}
	return *b.topLeftCellSize
}

func (b *Border) TopCellSize() Size {
	if b.topCellSize == nil {
		size := getStrSize(b.top)
		b.topCellSize = &size
	}
	return *b.topCellSize
}

func (b *Border) TopRightCellSize() Size {
	if b.topRightCellSize == nil {
		size := getStrSize(b.topRight)
		b.topRightCellSize = &size
	}
	return *b.topRightCellSize
}

func (b *Border) LeftCellSize() Size {
	if b.leftCellSize == nil {
		size := getStrSize(b.left)
		b.leftCellSize = &size
	}
	return *b.leftCellSize
}

func (b *Border) RightCellSize() Size {
	if b.rightCellSize == nil {
		size := getStrSize(b.right)
		b.rightCellSize = &size
	}
	return *b.rightCellSize
}

func (b *Border) BottomLeftCellSize() Size {
	if b.bottomLeftCellSize == nil {
		size := getStrSize(b.bottomLeft)
		b.bottomLeftCellSize = &size
	}
	return *b.bottomLeftCellSize
}

func (b *Border) BottomCellSize() Size {
	if b.bottomCellSize == nil {
		size := getStrSize(b.bottom)
		b.bottomCellSize = &size
	}
	return *b.bottomCellSize
}

func (b *Border) BottomRightCellSize() Size {
	if b.bottomRightCellSize == nil {
		size := getStrSize(b.bottomRight)
		b.bottomRightCellSize = &size
	}
	return *b.bottomRightCellSize
}

func (b *Border) LeftWidth() int {
	topLeftSize := b.TopLeftCellSize()
	leftSize := b.LeftCellSize()
	bottomLeftSize := b.BottomLeftCellSize()
	return max(topLeftSize.Width, leftSize.Width, bottomLeftSize.Width)
}

func (b *Border) RightWidth() int {
	topRightSize := b.TopRightCellSize()
	rightSize := b.RightCellSize()
	bottomRightSize := b.BottomRightCellSize()
	return max(topRightSize.Width, rightSize.Width, bottomRightSize.Width)
}

func (b *Border) TopHeight() int {
	topLeftSize := b.TopLeftCellSize()
	topSize := b.TopCellSize()
	topRightSize := b.TopRightCellSize()
	return max(topLeftSize.Height, topSize.Height, topRightSize.Height)
}

func (b *Border) BottomHeight() int {
	bottomLeftSize := b.BottomLeftCellSize()
	bottomSize := b.BottomCellSize()
	bottomRightSize := b.BottomRightCellSize()
	return max(bottomLeftSize.Height, bottomSize.Height, bottomRightSize.Height)
}

func getStrSize(str string) Size {
	_, info, err := ApplyNoWrap(false, Size{Width: math.MaxInt, Height: math.MaxInt}, str)
	if err != nil {
		panic(fmt.Errorf("Failed to calculate string size: %s", err))
	}
	return info.Size
}
