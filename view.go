package blitra

import (
	"errors"
	"os"
	"time"
)

const viewID = "__ROOT__"

type TargetBuffer int

const (
	PrimaryBuffer TargetBuffer = iota
	SecondaryBuffer
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

// Calling the View function returns a ViewHandle. ViewHandle provides
// methods for controlling the view, and rendering it.
//
// ViewHandle internally uses StdioManager to manage the target TTY.
type ViewHandle struct {
	fn            func(ViewState) any
	screenBuffer  *ScreenBuffer
	lastFrameTime time.Time

	opts   ViewOpts
	x      int
	y      int
	width  int
	height int
	state  ViewState

	stdioManager *StdioManager
}

// Given to the render function for the view, ViewState contains information
// about the state of the view, it's elements, and the frame being rendered.
type ViewState struct {
	elementIndex ElementIndex
	deltaTime    float64
	events       []Event
}

// Returns the delta time between the current frame and the previous frame.
func (v *ViewState) DeltaTime() float64 {
	return v.deltaTime
}

// Returns the size of the view.
//
// Note that this is distinct from the size of the tty. The view may be
// rendered into a smaller area of the terminal.
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

// Allows querying the size of an element by its ID.
//
// The size is able to be retrieved because it uses the already calculated
// size from the previous frame.
//
// WARNING: Because the size is from the previous frame, it will return a size
// of 0 on the first frame. Be ware of this when using the size for
// calculations.
//
// TODO:
//   - We should replace this method with one that returns a ElementHandle.
//   - ElementHandles should provide methods for querying the element's state,
//     including the size.
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

// Creates a ViewHandle with the given options and render function.
//
// The given render function will be called each frame to construct the view's
// internal element tree. Render functions can return any renderable type.
// Renderable types are:
// - string     - converted to a TextRenderable.
// - Renderable - any struct that implements the Renderable interface.
// - []any      - a list of renderables. It's of any so the list can be mixed.
// - nil        - nil can be used to skip rendering content.
func View(opts ViewOpts, fn func(ViewState) any) *ViewHandle {
	return &ViewHandle{
		opts:         opts,
		fn:           fn,
		stdioManager: NewStdioManager(opts.TTY),
	}
}

// Binds the view to the TTY.
//
// WARNING: It is possible to bind move than one view at a time, but views
// should not overlap. Overlapping views will produce undefined behavior.
func (v *ViewHandle) Bind() error {
	if err := v.stdioManager.Bind(); err != nil {
		return err
	}

	v.x = VOr(v.opts.X, 0)
	v.y = VOr(v.opts.Y, 0)
	v.width = VOr(v.opts.Width, v.stdioManager.ttySize.Width)
	v.height = VOr(v.opts.Height, v.stdioManager.ttySize.Height)
	v.screenBuffer = NewScreenBuffer(v.x, v.y, v.width, v.height)
	PrepareScreen(v)

	return nil
}

// Unbinds the view from the TTY, restoring the TTY to its previous state.
func (v *ViewHandle) Unbind() error {
	RestoreScreen(v)
	return v.stdioManager.Unbind()
}

// Should be called each frame, RenderFrame executes the view's render function,
// constructing an internal element tree. It then flows layout and renders the
// view.
func (v *ViewHandle) RenderFrame() ([]Event, error) {
	// Ensure that if a panic occurs while rendering the view, we
	// at least try restore the TTY to a usable state.
	defer func() {
		err := recover()
		if err != nil {
			if err2 := v.Unbind(); err2 != nil {
				panic(err2)
			}
			panic(err)
		}
	}()

	if !v.stdioManager.isBound {
		return nil, errors.New("view is not bound to a TTY. Make sure to call Bind before rendering")
	}

	events := v.stdioManager.TakeEvents()
	v.state.events = events

	frameTime := time.Now()
	if v.lastFrameTime.IsZero() {
		v.lastFrameTime = frameTime.Add(-time.Second / 60)
	}
	v.state.deltaTime = frameTime.Sub(v.lastFrameTime).Seconds()
	v.lastFrameTime = frameTime

	rootElement, elementIndex, err := ElementTreeAndIndexFromRenderable(&viewRenderable{view: v}, v.state)
	if err != nil || rootElement == nil {
		return nil, err
	}
	v.state.elementIndex = elementIndex

	if v.opts.Width == nil {
		v.width = v.stdioManager.ttySize.Width
	}
	if v.opts.Height == nil {
		v.height = v.stdioManager.ttySize.Height
	}
	rootElement.AvailableSize.Width = v.width
	rootElement.AvailableSize.Height = v.height
	rootElement.IntrinsicSize = rootElement.AvailableSize
	rootElement.Size = rootElement.AvailableSize

	v.screenBuffer.MaybeResize(v.x, v.y, v.width, v.height)

	if err := UpdateLayout(rootElement); err != nil {
		return nil, err
	}
	if err := Render(v, rootElement); err != nil {
		return nil, err
	}

	return events, nil
}

// This struct is used to wrap the render function of the view, so it implements
// the Renderable interface.
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
