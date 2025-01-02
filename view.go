package blitra

import (
	"errors"
	"os"
	"time"
)

const viewID = "__ROOT__"

// Options for controlling how a view is rendered.
type ViewOpts struct {
	// Axis controls the direction in which child elements are laid out.
	// The default is Vertically.
	Axis *Axis

	// Padding in columns from all edges of the view. Will be overridden by
	// LeftPadding, RightPadding, TopPadding, and/or BottomPadding.
	Padding *int
	// Padding in columns from the left edge of the view.
	LeftPadding *int
	// Padding in columns from the right edge of the view.
	RightPadding *int
	// Padding in rows from the top edge of the view.
	TopPadding *int
	// Padding in rows from the bottom edge of the view.
	BottomPadding *int
	// Sets the gap between child elements.
	Gap *int

	// Sets how child elements will be aligned along the axis. The default is
	// Stretch.
	Align *Align
	// Sets how child elements will be spaced along the axis within the view. The
	// default is Stretch.
	Justify *Justify

	// The X coordinates of the view. Useful if you want to render the
	// view into a specific location on the terminal rather than the whole.
	X *int
	// The Y coordinates of the view. Useful if you want to render the
	// view into a specific location on the terminal rather than the whole.
	Y *int

	// Sets a specific width in columns for the view.
	Width *int
	// Sets a specific height in rows for the view.
	Height *int

	// Sets the background color of the view. If unset the view will use the
	// default background color of the terminal.
	BackgroundColor *string
	// Sets the text color of the view. If unset the view will use the default
	// text color of the terminal.
	TextColor *string

	// The target TTY file to render into. Defaults to os.Stdout.
	TTY *os.File

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
	fn            func(ViewState) any
	tty           *os.File
	screenBuffer  *ScreenBuffer
	lastFrameTime time.Time

	x      int
	y      int
	width  int
	height int
	opts   ViewOpts
	state  ViewState
}

// Contains information about the current state of the view.
type ViewState struct {
	elementIndex ElementIndex
	deltaTime    float64
	// events       any
}

func (v *ViewState) DeltaTime() float64 {
	return v.deltaTime
}

func (v *ViewState) Size() Size {
	element, ok := v.elementIndex[viewID]
	if !ok {
		return Size{}
	}
	return Size{
		Width:  element.Size.Width,
		Height: element.Size.Height,
	}
}

func (v *ViewState) ElementSize(id string) Size {
	element, ok := v.elementIndex[id]
	if !ok {
		return Size{}
	}
	return Size{
		Width:  element.Size.Width - element.LeftMargin() - element.RightMargin(),
		Height: element.Size.Height - element.TopMargin() - element.BottomMargin(),
	}
}

// Creates a new view and returns a handle to it.
func View(opts ViewOpts, fn func(ViewState) any) *ViewHandle {
	return &ViewHandle{
		fn:   fn,
		opts: opts,
	}
}

// Binds the view to the terminal. This may involve switching to an alternate
// screen buffer, switching to raw mode, etc.
func (v *ViewHandle) Bind() error {
	if v.opts.TTY != nil {
		v.tty = v.opts.TTY
	} else {
		v.tty = os.Stdout
	}
	if !IsTTY(v.tty) {
		return errors.New("cannot bind. The target is not a TTY")
	}
	var termWidth, termHeight int
	if v.opts.Width == nil || v.opts.Height == nil {
		termWidth, termHeight = MustGetTerminalSize(v.tty)
	}
	v.x = V(v.opts.X)
	v.y = V(v.opts.Y)
	v.width = VOr(v.opts.Width, termWidth)
	v.height = VOr(v.opts.Height, termHeight)
	v.screenBuffer = NewScreenBuffer(v.x, v.y, v.width, v.height)

	PrepareScreen(v)

	return nil
}

// Unbinds the view from the terminal restoring the terminal to its previous
// state.
func (v *ViewHandle) Unbind() {
	// MustRestoreTTYToNormal(v.tty, v.terminalState)

	RestoreScreen(v)
}

// Renders a frame based on the current state of the view.
func (v *ViewHandle) RenderFrame() error {
	if v.tty == nil {
		panic("cannot render. The view is not bound")
	}

	frameTime := time.Now()
	if v.lastFrameTime.IsZero() {
		v.lastFrameTime = frameTime.Add(-time.Second / 60)
	}
	v.state.deltaTime = frameTime.Sub(v.lastFrameTime).Seconds()
	v.lastFrameTime = frameTime

	rootElement, elementIndex, err := ElementTreeAndIndexFromRenderable(&viewRenderable{view: v}, v.state)
	if err != nil || rootElement == nil {
		return err
	}
	v.state.elementIndex = elementIndex
	rootElement.AvailableSize.Width = v.width
	rootElement.AvailableSize.Height = v.height
	rootElement.IntrinsicSize = rootElement.AvailableSize
	rootElement.Size = rootElement.AvailableSize

	UpdateLayout(rootElement)
	Render(v, rootElement)

	return nil
}

type viewRenderable struct {
	view *ViewHandle
}

var _ Renderable = &viewRenderable{}

func (v *viewRenderable) ID() string {
	return viewID
}

func (v *viewRenderable) Style() Style {
	return Style{
		Axis:            v.view.opts.Axis,
		LeftPadding:     OrP(v.view.opts.LeftPadding, v.view.opts.Padding),
		RightPadding:    OrP(v.view.opts.RightPadding, v.view.opts.Padding),
		TopPadding:      OrP(v.view.opts.TopPadding, v.view.opts.Padding),
		BottomPadding:   OrP(v.view.opts.BottomPadding, v.view.opts.Padding),
		Gap:             v.view.opts.Gap,
		Align:           v.view.opts.Align,
		Justify:         v.view.opts.Justify,
		BackgroundColor: v.view.opts.BackgroundColor,
		TextColor:       v.view.opts.TextColor,
	}
}

func (v *viewRenderable) Render(state ViewState) any {
	return v.view.fn(state)
}
