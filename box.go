package blitra

type BoxRenderable struct {
	id   string
	opts BoxOpts
	fn   func(ctx BoxState) any
}

type BoxOpts struct {

	// Can be set to vertical or horizontal with a default of horizontal.
	// Determines the axis that the children elements will be laid out on.
	Axis *Axis

	Grow *int

	Shrink *int

	// The alignment of the box's children. Defaults to stretch.
	Align *Align
	// The justification of the box's children. Defaults to stretch.
	Justify *Justify

	// How many empty columns to the left of the box's children.
	LeftPadding *int
	// How many empty columns to the right of the box's children.
	RightPadding *int
	// How many empty rows above the box's children.
	TopPadding *int
	// How many empty rows below the box's children.
	BottomPadding *int
	// How many empty columns and rows around the box's children. Overridden by
	// the other padding values.
	Padding *int
	// How many empty columns/rows between each child element.
	Gap *int

	// How many empty columns to the left of the box.
	LeftMargin *int
	// How many empty columns to the right of the box.
	RightMargin *int
	// How many empty rows above the box.
	TopMargin *int
	// How many empty rows below the box.
	BottomMargin *int
	// How many empty columns and rows around the box. Overridden by the other
	// margin values.
	Margin *int

	// The width of the box in columns. If 0 the width will be determined
	// automatically.
	Width *int
	// The minimum width of the box in columns. If 0 the minimum width will be
	// determined automatically.
	MinWidth *int
	// The maximum width of the box in columns. If 0 the maximum width will be
	// determined automatically.
	MaxWidth *int

	// The height of the box in rows. If 0 the height will be determined
	// automatically.
	Height *int
	// The minimum height of the box in rows. If 0 the minimum height will be
	// determined automatically.
	MinHeight *int
	// The maximum height of the box in rows. If 0 the maximum height will be
	// determined automatically.
	MaxHeight *int

	// The border style of the top of the box.
	LeftBorder *Border
	// The border style of the right of the box.
	RightBorder *Border
	// The border style of the top of the box.
	TopBorder *Border
	// The border style of the bottom of the box.
	BottomBorder *Border
	// The border style of the box. Overridden by the other border values.
	Border *Border

	// How text should wrap in the box. Defaults to WordWrap.
	TextWrap *TextWrap
	// If true, when text cannot fit in the box it will be truncated with an
	// ellipsis.
	Ellipsis *bool

	// The background color of the box.
	BackgroundColor *string
	// The text color of the box. This will be inherited by the box's children.
	TextColor *string

	DEBUG_ID string
}

var _ Renderable = &BoxRenderable{}

// Allows dividing views into horizontal or vertical sections. Also provides
// layout options as to control spacing and alignment.
func Box(id string, opts BoxOpts, fn func(ctx BoxState) any) *BoxRenderable {
	return &BoxRenderable{
		id:   id,
		opts: opts,
		fn:   fn,
	}
}

func (b *BoxRenderable) ID() string {
	return b.id
}

func (b *BoxRenderable) Style() Style {
	return Style{
		DEBUG_ID: b.opts.DEBUG_ID,

		Grow: b.opts.Grow,
		Axis: b.opts.Axis,

		LeftPadding:   OrP(b.opts.LeftPadding, b.opts.Padding),
		RightPadding:  OrP(b.opts.RightPadding, b.opts.Padding),
		TopPadding:    OrP(b.opts.TopPadding, b.opts.Padding),
		BottomPadding: OrP(b.opts.BottomPadding, b.opts.Padding),

		Gap: b.opts.Gap,

		LeftMargin:   OrP(b.opts.LeftMargin, b.opts.Margin),
		RightMargin:  OrP(b.opts.RightMargin, b.opts.Margin),
		TopMargin:    OrP(b.opts.TopMargin, b.opts.Margin),
		BottomMargin: OrP(b.opts.BottomMargin, b.opts.Margin),

		Width:    b.opts.Width,
		MinWidth: b.opts.MinWidth,
		MaxWidth: b.opts.MaxWidth,

		Height:    b.opts.Height,
		MinHeight: b.opts.MinHeight,
		MaxHeight: b.opts.MaxHeight,

		Align:   b.opts.Align,
		Justify: b.opts.Justify,

		LeftBorder:   OrP(b.opts.LeftBorder, b.opts.Border),
		RightBorder:  OrP(b.opts.RightBorder, b.opts.Border),
		TopBorder:    OrP(b.opts.TopBorder, b.opts.Border),
		BottomBorder: OrP(b.opts.BottomBorder, b.opts.Border),

		TextWrap: b.opts.TextWrap,
		Ellipsis: b.opts.Ellipsis,

		BackgroundColor: b.opts.BackgroundColor,
		TextColor:       b.opts.TextColor,
	}
}

// Implements the Renderable interface.
func (b *BoxRenderable) Render(state ViewState) any {
	if b.fn == nil {
		return nil
	}
	boxState := BoxState{}
	return b.fn(boxState)
}

type BoxState struct {
	// Indicates if the view was clicked.
	Clicked bool
}
