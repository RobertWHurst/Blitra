package format

import "fmt"

func Reset() {
	fmt.Print("\x1B[0m")
}

func Bold() {
	fmt.Print("\x1B[1m")
}

func UnBold() {
	fmt.Print("\x1B[22m")
}

func Dim() {
	fmt.Print("\x1B[2m")
}

func Italic() {
	fmt.Print("\x1B[3m")
}

func UnItalic() {
	fmt.Print("\x1B[23m")
}

func Underline() {
	fmt.Print("\x1B[4m")
}

func UnUnderline() {
	fmt.Print("\x1B[24m")
}

func Blink() {
	fmt.Print("\x1B[5m")
}

func UnBlink() {
	fmt.Print("\x1B[25m")
}

func BlinkRapid() {
	fmt.Print("\x1B[6m")
}

func Reverse() {
	fmt.Print("\x1B[7m")
}

func UnReverse() {
	fmt.Print("\x1B[27m")
}

func Hidden() {
	fmt.Print("\x1B[8m")
}

func UnHidden() {
	fmt.Print("\x1B[28m")
}

func Strike() {
	fmt.Print("\x1B[9m")
}

func UnStrike() {
	fmt.Print("\x1B[29m")
}
