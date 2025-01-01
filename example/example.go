package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	b "github.com/RobertWHurst/blitra"
)

func main() {

	view := b.View(b.ViewOpts{
		Axis:            b.P(b.VerticalAxis),
		Align:           b.P(b.CenterAlign),
		Justify:         b.P(b.CenterJustify),
		Gap:             b.P(2),
		BackgroundColor: b.P("#002"),
		TargetBuffer:    b.SecondaryBuffer,
	}, func(ctx b.ViewState) any {

		return []any{
			b.Box(b.BoxOpts{
				DEBUG_ID:        "box",
				Align:           b.P(b.CenterAlign),
				TopPadding:      b.P(1),
				BottomPadding:   b.P(1),
				LeftPadding:     b.P(2),
				RightPadding:    b.P(2),
				Gap:             b.P(2),
				Border:          b.DoubleBorder(),
				TextColor:       b.P("#005"),
				BackgroundColor: b.P("#eef"),
			}, func(_ b.BoxState) any {
				now := time.Now()

				return []any{

					b.Box(b.BoxOpts{
						DEBUG_ID: "sub-box-1",
					}, func(_ b.BoxState) any {
						return "The time is:"
					}),

					b.Box(b.BoxOpts{
						DEBUG_ID:        "sub-box-2",
						Border:          b.DoubleBorder(),
						BackgroundColor: b.P("#ccd"),
					}, func(_ b.BoxState) any {
						if ctx.Clicked {
							return "Hello, mouse click!"
						}
						return now.Format("2006-01-02 15:04:05.000")
					}),
				}
			}),

			b.Box(b.BoxOpts{
				DEBUG_ID:  "box-2",
				TextColor: b.P("#f00"),
				Width:     b.P(10),
			}, func(_ b.BoxState) any {
				return "Delta time is: " + strconv.FormatFloat(ctx.DeltaTime, 'f', 6, 64)
			}),
		}
	})

	if err := view.Bind(); err != nil {
		panic(err)
	}

	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt)

	var frameCount int
	var frameTime time.Duration
loop:
	for {
		select {
		case <-osSignalChan:
			break loop
		default:
		}
		start := time.Now()
		view.RenderFrame()
		elapsed := time.Since(start)
		frameTime += elapsed
		frameCount += 1
		time.Sleep(time.Second / 60)
	}

	view.Unbind()

	fmt.Println("Average frame time:", frameTime/time.Duration(frameCount))
}
