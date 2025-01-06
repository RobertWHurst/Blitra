package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RobertWHurst/blitra"
)

func main() {
	clockView := createClockView()
	if err := clockView.Bind(); err != nil {
		panic(err)
	}
	defer func() {
		if err := clockView.Unbind(); err != nil {
			panic(err)
		}
	}()

	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)

loop:
	for {
		select {
		case <-osSignalChan:
			break loop
		default:
		}

		events, err := clockView.RenderFrame()
		if err != nil {
			panic(err)
		}

		DebugLogEvents(events)

		for _, event := range events {
			if event.Kind == blitra.CtrlKeyEvent && event.ModifiedChar == 'C' {
				break loop
			}
		}

		// 120 FPS
		time.Sleep(time.Second / 120)
	}

	fmt.Println("Goodbye!")
}
