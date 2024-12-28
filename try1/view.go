package tui

// Options for controlling how a view is rendered.
type ViewOpts struct {
	// Axis controls the direction in which child elements are laid out.
	// The default is Vertically.
	Axis Axis

	// Padding in columns from all edges of the view. Will be overridden by
	// PaddingLeft, PaddingRight, PaddingTop, and/or PaddingBottom.
	Padding int
	// Padding in columns from the left edge of the view.
	PaddingLeft int
	// Padding in columns from the right edge of the view.
	PaddingRight int
	// Padding in rows from the top edge of the view.
	PaddingTop int
	// Padding in rows from the bottom edge of the view.
	PaddingBottom int

	// Sets how child elements will be aligned along the axis. The default is
	// Stretch.
	Align Align

	// Sets how child elements will be spaced along the axis within the view. The
	// default is Stretch.
	Justify Justify

	// Sets the background color of the view. If unset the view will use the
	// default background color of the terminal.
	BackgroundColor string

	// Sets the text color of the view. If unset the view will use the default
	// text color of the terminal.
	TextColor string
}

// Contains information about the current state of the view.
type ViewContext struct {
	// Indicates if the view is clicked.
	Clicked bool
	// A slice of keys that are currently pressed.
	KeysPressed []string
}

// Has methods for rendering the view frame by frame.
type ViewHandle struct {
	opts       ViewOpts
	renderFn   func(ctx ViewContext) any
	frameChars [][]rune
}

// Creates a new view and returns a handle to it.
func View(opts ViewOpts, renderFn func(ctx ViewContext) any) *ViewHandle {
	return &ViewHandle{
		opts:     opts,
		renderFn: renderFn,
	}
}

// Options for controlling how a view is rendered for a single frame.
type RenderOpts struct {
	// The position of the mouse click in the view. The first element is the x
	// position in columns and the second element is the y position in rows.
	// should be nil if no click.
	ClickAt [2]int
	// A slice of keys that are currently pressed.
	KeysPressed []string
	// The width of the view in columns. Likely should be the column count of the
	// terminal.
	Width int
	// The height of the view in rows. Likely should be the row count of the
	// terminal.
	Height int
}

// string -wrapWithTextRenderable-> []Renderable
// func -call(CustomCtx)-> []Renderable
// Renderable -Render(RenderCtx)-> []Renderable

// Renders the view for a single frame.
func (v *ViewHandle) RenderFrame(opts RenderOpts) {
	viewCtx := ViewContext{}
	renderCtx := RenderContext{}

	renderables := v.renderFn(viewCtx)

	rootElement := ViewElement{}
}

type VO = ViewOpts
type VC = ViewContext

var V = View
