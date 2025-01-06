package main

import (
	"fmt"
	"os"

	"github.com/RobertWHurst/blitra"
)

func DebugLogEvents(events []blitra.Event) {
	debugFile, err := os.OpenFile("events-debug-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		switch event.Kind {
		case blitra.CtrlKeyEvent:
			if event.Key == blitra.NoKey {
				fmt.Fprintf(debugFile, "CtrlKeyEvent: %c\n", event.ModifiedChar)
			} else {
				fmt.Fprintf(debugFile, "CtrlKeyEvent: %s\n", event.Key)
			}
		case blitra.AltKeyEvent:
			if event.Key == blitra.NoKey {
				fmt.Fprintf(debugFile, "AltKeyEvent: %c\n", event.ModifiedChar)
			} else {
				fmt.Fprintf(debugFile, "AltKeyEvent: %s\n", event.Key)
			}
		case blitra.ShiftKeyEvent:
			fmt.Fprintf(debugFile, "ShiftKeyEvent: %s\n", event.Key)
		case blitra.KeyEvent:
			fmt.Fprintf(debugFile, "KeyEvent: %s\n", event.Key)
		case blitra.FocusEvent:
			fmt.Fprintln(debugFile, "FocusEvent")
		case blitra.BlurEvent:
			fmt.Fprintln(debugFile, "BlurEvent")
		case blitra.MouseDownEvent:
			fmt.Fprintf(debugFile, "MouseDownEvent: button=%d, x=%d, y=%d\n", event.MouseButton, event.MouseX, event.MouseY)
		case blitra.MouseUpEvent:
			fmt.Fprintf(debugFile, "MouseUpEvent: button=%d, x=%d, y=%d\n", event.MouseButton, event.MouseX, event.MouseY)
		case blitra.MouseMoveEvent:
			if event.MouseButton != blitra.NoMouseButton {
				fmt.Fprintf(debugFile, "MouseMoveEvent: button=%d, x=%d, y=%d\n", event.MouseButton, event.MouseX, event.MouseY)
			} else {
				fmt.Fprintf(debugFile, "MouseMoveEvent: x=%d, y=%d\n", event.MouseX, event.MouseY)
			}
		case blitra.MouseScrollEvent:
			fmt.Fprintf(debugFile, "MouseScrollEvent: direction=%d\n", event.MouseScrollDirection)
		case blitra.CharInputEvent:
			fmt.Fprintf(debugFile, "CharInputEvent: %c\n", event.Char)
		}
	}
}
