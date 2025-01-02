package main

import (
	"strconv"
	"time"

	b "github.com/RobertWHurst/blitra"
)

func createClockView() *b.ViewHandle {

	// Set some variables the view will use to animate the clock position.
	offsetY := 0.
	offsetX := 0.
	offsetXSize := 1.2
	offsetYSize := 0.8

	// Create a view that will render the clock.
	return b.View(viewStyle, func(view b.ViewState) any {
		// This function will be called every frame to render the view.

		// Get the current time.
		now := time.Now()

		// Update the clock position and apply it to the content box style.
		updateOffsets(view, &offsetX, &offsetY, &offsetXSize, &offsetYSize)
		appStyleWPos := appStyle
		appStyleWPos.LeftMargin = b.P(int(offsetX))
		appStyleWPos.TopMargin = b.P(int(offsetY))

		// Setup the box structure for the view.
		return b.Box("content", appStyleWPos, func(_ b.BoxState) any {
			return []any{
				b.Box("clock", clockStyle, func(clock b.BoxState) any {
					return []any{
						b.Box("label", labelStyle, func(_ b.BoxState) any {
							return "The time is:"
						}),
						b.Box("value", valueStyle, func(_ b.BoxState) any {
							if clock.Clicked {
								return "Hello, mouse click!"
							}
							return now.Format("2006-01-02 15:04:05.000")
						}),
					}
				}),
				b.Box("deltaTimeDebug", deltaTimeDebugStyle, func(_ b.BoxState) any {
					return "Delta time is: " + strconv.FormatFloat(view.DeltaTime(), 'f', 6, 64)
				}),
			}
		})
	})
}

var viewStyle = b.ViewOpts{
	BackgroundColor: b.P("#002"),
	TargetBuffer:    b.SecondaryBuffer,
	Align:           b.P(b.StartAlign),
}

var appStyle = b.BoxOpts{
	Axis:  b.P(b.VerticalAxis),
	Gap:   b.P(2),
	Align: b.P(b.CenterAlign),
}

var clockStyle = b.BoxOpts{
	Align:           b.P(b.CenterAlign),
	TopPadding:      b.P(1),
	BottomPadding:   b.P(1),
	LeftPadding:     b.P(2),
	RightPadding:    b.P(2),
	Gap:             b.P(2),
	Border:          b.DoubleBorder(),
	TextColor:       b.P("#005"),
	BackgroundColor: b.P("#eef"),
}

var valueStyle = b.BoxOpts{
	Border:          b.RoundBorder(),
	TextColor:       b.P("#fff"),
	BackgroundColor: b.P("#007"),
}

var labelStyle = b.BoxOpts{}

var deltaTimeDebugStyle = b.BoxOpts{
	TextColor: b.P("#a33"),
}

func updateOffsets(view b.ViewState, offsetX, offsetY, offsetXSize, offsetYSize *float64) {
	deltaTime := view.DeltaTime()

	*offsetX += *offsetXSize * deltaTime
	*offsetY += *offsetYSize * deltaTime

	viewSize := view.Size()
	rootBoxSize := view.ElementSize("content")

	maxOffsetX := float64(viewSize.Width - rootBoxSize.Width)
	maxOffsetY := float64(viewSize.Height - rootBoxSize.Height)

	if *offsetX > maxOffsetX {
		*offsetX = maxOffsetX
		*offsetXSize = -*offsetXSize
	}
	if *offsetX < 0 {
		*offsetX = 0
		*offsetXSize = -*offsetXSize
	}
	if *offsetY > maxOffsetY {
		*offsetY = maxOffsetY
		*offsetYSize = -*offsetYSize
	}
	if *offsetY < 0 {
		*offsetY = 0
		*offsetYSize = -*offsetYSize
	}
}
