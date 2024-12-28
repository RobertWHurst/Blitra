package tui

type DivisionOpts struct {
	Grow bool
	Axis Axis

	PaddingLeft   int
	PaddingRight  int
	PaddingTop    int
	PaddingBottom int
	Padding       int

	MarginLeft   int
	MarginRight  int
	MarginTop    int
	MarginBottom int
	Margin       int

	Width     int
	MinWidth  int
	MaxWidth  int
	Height    int
	MinHeight int
	MaxHeight int

	Align   Align
	Justify Justify

	BackgroundColor string
	TextColor       string

	Border Border
}

type DivisionContext struct {
	Clicked     bool
	KeysPressed []string
}

type DivisionElement struct {
	opts     DivisionOpts
	RenderFn func(DivisionContext) any
}

func Division(opts DivisionOpts, renderFn func(ctx DivisionContext) any) DivisionElement {
	return DivisionElement{
		opts:     opts,
		RenderFn: renderFn,
	}
}

func (d DivisionElement) Render(ctx RenderContext) any {
	dCtx := DivisionContext{}
	d.RenderFn(dCtx)
}

type DO = DivisionOpts

var D = Division
