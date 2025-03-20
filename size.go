package blitra

import "math"

type Size struct {
	Width  int
	Height int
}

var MaxSize = Size{Width: math.MaxInt, Height: math.MaxInt}
