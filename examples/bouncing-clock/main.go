package main

import (
	"fmt"
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
		clockView.RenderFrame()
		elapsed := time.Since(start)
		frameTime += elapsed
		frameCount += 1
		time.Sleep(time.Second / 60)
	}

	fmt.Println("Average frame time:", frameTime/time.Duration(frameCount))
}
