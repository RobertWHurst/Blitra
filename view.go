package blitra

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

	// Sets the target buffer to render the view into. If unset the view will
	// render into the terminal.
	TargetBuffer TargetBuffer

	// By default the view will intercept stdout so if print logging is done it
	// will be saved until the view is unbound, after which it will be printed.
	// This can also be seen if debugging is enabled in the debug mode for the
	// view. If set to true this behavior will be disabled. The view will
	// still render, but could be corrupted by print logging interfering with
	// the view's output. Leaving this as false is recommended unless you have
	// a good reason disable interception.
	DisableStdoutInterception bool
}

// Can be used to control rendering of the view.
type ViewHandle struct {
	opts  ViewOpts
	state ViewState
	fn    func(ViewState) any
}

// Creates a new view and returns a handle to it.
func View(opts ViewOpts, fn func(ViewState) any) ViewHandle {
	return ViewHandle{
		opts: opts,
		fn:   fn,
	}
}

// Binds the view to the terminal. This may involve switching to an alternate
// screen buffer, switching to raw mode, etc.
func (v *ViewHandle) Bind() {
}

// Unbinds the view from the terminal restoring the terminal to its previous
// state.
func (v *ViewHandle) Unbind() {
}

// Renders a frame based on the current state of the view.
func (v *ViewHandle) RenderFrame() {
	_ = ElementFromAnyWithState(v.fn(v.state), v.state)
}

// Contains information about the current state of the view.
type ViewState struct {
	// Indicates at which column and row the view was clicked.
	ClickAt Point
	// Indicates if the view was clicked.
	Clicked bool
	// Indicates which keys are currently pressed.
	KeysPressed []Key
}
