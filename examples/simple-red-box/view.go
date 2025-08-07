package main

import (
	b "github.com/RobertWHurst/blitra"
)

func createRedBoxView() *b.ViewHandle {
	return b.View(viewStyle, func(view b.ViewState) any {
		return b.Box("red-box", redBoxStyle, func(_ b.BoxState) any {
			return "Lorem ipsum odor amet, consectetuer adipiscing elit. Quam a purus per lectus amet eros cras; elementum egestas. Curae accumsan conubia, quisque vulputate nascetur maecenas. Quis cras sollicitudin himenaeos lobortis venenatis torquent nibh bibendum."
		})
	})
}

var viewStyle = b.ViewOpts{
	BackgroundColor: b.P("#fff"),
	TargetBuffer:    b.SecondaryBuffer,
	Align:           b.P(b.StartAlign),
	Axis:            b.P(b.VerticalAxis),
}

var redBoxStyle = b.BoxOpts{
	TopMargin:       b.P(1),
	LeftMargin:      b.P(2),
	Width:           b.P(20),
	Height:          b.P(10),
	BackgroundColor: b.P("#f00"),
	TextColor:       b.P("#fff"),
	Ellipsis:        b.P(true),
}
