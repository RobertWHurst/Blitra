package color

import (
	"strconv"

	"github.com/robertwhurst/blitra/stdio"
)

type ColorMode int

const (
	Mode256 ColorMode = iota
	Mode16
)

// Standard 16 color palette
//
// Foreground colors
func Black() {
	stdio.Print("\x1B[30m")
}

func Red() {
	stdio.Print("\x1B[31m")
}

func Green() {
	stdio.Print("\x1B[32m")
}

func Yellow() {
	stdio.Print("\x1B[33m")
}

func Blue() {
	stdio.Print("\x1B[34m")
}

func Magenta() {
	stdio.Print("\x1B[35m")
}

func Cyan() {
	stdio.Print("\x1B[36m")
}

func White() {
	stdio.Print("\x1B[37m")
}

// background colors
func BlackBackground() {
	stdio.Print("\x1B[40m")
}

func RedBackground() {
	stdio.Print("\x1B[41m")
}

func GreenBackground() {
	stdio.Print("\x1B[42m")
}

func YellowBackground() {
	stdio.Print("\x1B[43m")
}

func BlueBackground() {
	stdio.Print("\x1B[44m")
}

func MagentaBackground() {
	stdio.Print("\x1B[45m")
}

func CyanBackground() {
	stdio.Print("\x1B[46m")
}

func WhiteBackground() {
	stdio.Print("\x1B[47m")
}

// Bright foreground colors

func BrightBlack() {
	stdio.Print("\x1B[90m")
}

func BrightRed() {
	stdio.Print("\x1B[91m")
}

func BrightGreen() {
	stdio.Print("\x1B[92m")
}

func BrightYellow() {
	stdio.Print("\x1B[93m")
}

func BrightBlue() {
	stdio.Print("\x1B[94m")
}

func BrightMagenta() {
	stdio.Print("\x1B[95m")
}

func BrightCyan() {
	stdio.Print("\x1B[96m")
}

func BrightWhite() {
	stdio.Print("\x1B[97m")
}

// Bright background colors

func BrightBlackBackground() {
	stdio.Print("\x1B[100m")
}

func BrightRedBackground() {
	stdio.Print("\x1B[101m")
}

func BrightGreenBackground() {
	stdio.Print("\x1B[102m")
}

func BrightYellowBackground() {
	stdio.Print("\x1B[103m")
}

func BrightBlueBackground() {
	stdio.Print("\x1B[104m")
}

func BrightMagentaBackground() {
	stdio.Print("\x1B[105m")
}

func BrightCyanBackground() {
	stdio.Print("\x1B[106m")
}

func BrightWhiteBackground() {
	stdio.Print("\x1B[107m")
}

// 256 colors and RGB
func Color256(color int) {
	if color < 0 || color > 255 {
		panic("Invalid color")
	}
	stdio.Printf("\x1B[38;5;%dm", color)
}

func Color256Backgound(color int) {
	if color < 0 || color > 255 {
		panic("Invalid color")
	}
	stdio.Printf("\x1B[48;5;%dm", color)
}

func ColorRGB(r, g, b int) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		panic("Invalid color")
	}
	stdio.Printf("\x1B[38;2;%d;%d;%dm", r, g, b)
}

func ColorRGBBackground(r, g, b int) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		panic("Invalid color")
	}
	stdio.Printf("\x1B[48;2;%d;%d;%dm", r, g, b)
}

func ColorHEX(hex string) {
	r, g, b := hexToRGB(hex)
	ColorRGB(int(r), int(g), int(b))
}

func ColorHEXBackground(hex string) {
	r, g, b := hexToRGB(hex)
	ColorRGBBackground(int(r), int(g), int(b))
}

func hexToRGB(hex string) (int, int, int) {
	if len(hex) != 6 || len(hex) != 3 {
		panic("Invalid hex color")
	}
	if len(hex) == 3 {
		hex = hex[0:1] + hex[0:1] + hex[1:2] + hex[1:2] + hex[2:3] + hex[2:3]
	}
	r, err := strconv.ParseInt(hex[0:2], 16, 0)
	if err != nil {
		panic(err)
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 0)
	if err != nil {
		panic(err)
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 0)
	if err != nil {
		panic(err)
	}
	return int(r), int(g), int(b)
}
