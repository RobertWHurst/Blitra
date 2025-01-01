package blitra

import (
	"errors"
	"os"
	"time"

	"golang.org/x/term"
)

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
	realStdout    *os.File
	stdoutChan    chan string
	screenBuffer  *ScreenBuffer
	lastFrameTime time.Time

	x             int
	y             int
	width         int
	height        int
	opts          ViewOpts
	state         ViewState
	terminalState *term.State
}

// Contains information about the current state of the view.
type ViewState struct {
	// The delta time between the last frame and the one before it.
	DeltaTime float64
	// Indicates at which column and row the view was clicked.
	ClickAt Point
	// Indicates if the view was clicked.
	Clicked bool
	// Indicates which keys are currently pressed.
	KeysPressed []Key
}

// Creates a new view and returns a handle to it.
func View(opts ViewOpts, fn func(ViewState) any) ViewHandle {
	return ViewHandle{
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

	// v.terminalState = MustSwitchTTYToRaw(v.tty)

	return nil
}

// Unbinds the view from the terminal restoring the terminal to its previous
// state.
func (v *ViewHandle) Unbind() {
	// MustRestoreTTYToNormal(v.tty, v.terminalState)

	RestoreScreen(v)
}

// Renders a frame based on the current state of the view.
func (v *ViewHandle) RenderFrame() {
	if v.tty == nil {
		panic("cannot render. The view is not bound")
	}

	frameTime := time.Now()
	if v.lastFrameTime.IsZero() {
		v.lastFrameTime = frameTime.Add(-time.Second / 60)
	}
	v.state.DeltaTime = frameTime.Sub(v.lastFrameTime).Seconds()
	v.lastFrameTime = frameTime

	rootElement := ElementFromRenderable(viewRenderable{
		view: v,
	}, v.state)
	if rootElement == nil {
		return
	}
	rootElement.AvailableSize.Width = v.width
	rootElement.AvailableSize.Height = v.height

	UpdateLayout(rootElement)
	Render(v, rootElement)

}

type viewRenderable struct {
	view *ViewHandle
}

var _ Renderable = viewRenderable{}

func (v viewRenderable) Style() Style {
	return Style{
		Axis:            v.view.opts.Axis,
		LeftPadding:     OrP(v.view.opts.LeftPadding, v.view.opts.Padding),
		RightPadding:    OrP(v.view.opts.RightPadding, v.view.opts.Padding),
		TopPadding:      OrP(v.view.opts.TopPadding, v.view.opts.Padding),
		BottomPadding:   OrP(v.view.opts.BottomPadding, v.view.opts.Padding),
		Align:           v.view.opts.Align,
		Justify:         v.view.opts.Justify,
		BackgroundColor: v.view.opts.BackgroundColor,
		TextColor:       v.view.opts.TextColor,
	}
}

func (v viewRenderable) Render(state ViewState) any {
	return v.view.fn(state)
}
