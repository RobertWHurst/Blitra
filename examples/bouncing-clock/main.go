package main

import (
	"os"
	"os/signal"
	"time"
)

func main() {
	clockView := createClockView()
	if err := clockView.Bind(); err != nil {
		panic(err)
	}
	defer clockView.Unbind()

	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt)

loop:
	for {
		select {
		case <-osSignalChan:
			break loop
		default:
		}

		clockView.RenderFrame()

		// 120 FPS
		time.Sleep(time.Second / 120)
	}
}
