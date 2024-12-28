package tui

import "fmt"

type RenderContext struct {
	// frameChars is a slice of a slice of runes. This slice will be used to
	// render the parent view's frame.
	frameChars [][]rune

	// These indicate the dimensions of the context
	col    int
	row    int
	width  int
	height int

	ClickAt     [2]int
	KeysPressed []string

	MaxWidth  int
	MaxHeight int

	TextColor       string
	BackgroundColor string
}

// Creates a sub context meant for rendering child elements. The col and row
// values are relative to the parent context.
func (v *RenderContext) SubContext(col, row, width, height int) RenderContext {
	return RenderContext{
		frameChars:      v.frameChars,
		col:             v.col + col,
		row:             v.col + row,
		width:           width,
		height:          height,
		ClickAt:         v.ClickAt,
		KeysPressed:     v.KeysPressed,
		MaxWidth:        v.MaxWidth,
		MaxHeight:       v.MaxHeight,
		TextColor:       v.TextColor,
		BackgroundColor: v.BackgroundColor,
	}
}

func (v *RenderContext) SubRender(result any) {
	elements := []Renderable{}
	switch result := result.(type) {
	case []Renderable:
		elements = result
	case Renderable:
		elements = []Renderable{result}
	case func() []Renderable:
		elements = result()
	case func() Renderable:
		elements = []Renderable{result()}
	case string:
		elements = []Renderable{Text(result)}
	default:
		if result != nil {
			panic("Cannot render type " + fmt.Sprintf("%T", result))
		}
	}
}

func (v *RenderContext) Draw(str string) {
	strChars := []rune(str)
	col := v.col
	row := v.row
	for _, char := range strChars {
		v.frameChars[row][col] = char
		col += 1
		if col == v.width {
			row += 1
			col = 0
		}
		if row == v.height {
			// We cannot draw past the last row
			break
		}
	}
}
