package main

import (
	"time"

	b "github.com/RobertWHurst/blitra"
)

func main() {

	view := b.View(b.ViewOpts{
		Align:           b.P(b.CenterAlign),
		Justify:         b.P(b.CenterJustify),
		BackgroundColor: b.P("#fef"),
		TargetBuffer:    b.SecondaryBuffer,
	}, func(ctx b.ViewState) any {

		return b.Box(b.BoxOpts{
			DEBUG_ID:        "box",
			Align:           b.P(b.CenterAlign),
			Gap:             b.P(1),
			Border:          b.DoubleBorder(),
			TextColor:       b.P("#005"),
			BackgroundColor: b.P("#afa"),
		}, func(_ b.BoxState) any {
			return []any{

				b.Box(b.BoxOpts{
					DEBUG_ID: "sub-box-1",
				}, func(_ b.BoxState) any {
					return "Example APP:"
				}),

				b.Box(b.BoxOpts{
					DEBUG_ID:        "sub-box-2",
					Border:          b.DoubleBorder(),
					BackgroundColor: b.P("#ffa"),
				}, func(_ b.BoxState) any {
					if ctx.Clicked {
						return "Hello, mouse click!"
					}
					return "Hello, world!"
				}),
			}
		})

	})

	if err := view.Bind(); err != nil {
		panic(err)
	}

	view.RenderFrame()
	time.Sleep(5 * time.Second)
	view.Unbind()
}
