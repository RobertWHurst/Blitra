package main

import "github.com/robertwhurst/blitra"

func main() {
	blitra.View(blitra.ViewOpts{}, func(ctx blitra.ViewContext) any {
		if ctx.Clicked {
			return "You clicked!"
		}

		return blitra.Division(blitra.DivisionOpts{
			Align:   blitra.AlignCenter,
			Justify: blitra.JustifyCenter,
		}, "Hello, world!")

	}).Render(blitra.RenderOpts{
		ClickAt: [2]int{0, 0},
	})
}

func MyTextComponent(text string) string {
	return text
}
