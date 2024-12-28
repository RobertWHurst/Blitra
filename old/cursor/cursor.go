package cursor

import (
	"fmt"

	"github.com/robertwhurst/blitra/stdio"
	"github.com/robertwhurst/blitra/terminal"
)

func Home() {
	stdio.Print("\x1B[H")
}

func Up(distance int) {
	stdio.Printf("\x1B[%dA", distance)
}

func UpWithScroll() {
	stdio.Printf("\x1BM")
}

func Down(distance int) {
	stdio.Printf("\x1B[%dB", distance)
}

func Right(distance int) {
	stdio.Printf("\x1B[%dC", distance)
}

func Left(distance int) {
	stdio.Printf("\x1B[%dD", distance)
}

func NextLine() {
	stdio.Printf("\x1B[E")
}

func PreviousLine() {
	stdio.Printf("\x1B[F")
}

func Column(column int) {
	stdio.Printf("\x1B[%dG", column)
}

func Position(column, row int) {
	stdio.Printf("\x1B[%d;%dH", row, column)
}

func GetPosition() (int, int, error) {
	if !terminal.Compatible() {
		return 0, 0, fmt.Errorf("stdout is not a terminal")
	}

	h, err := terminal.RawMode()
	if err != nil {
		return 0, 0, err
	}
	defer h.Restore()

	stdio.Print("\x1B[6n")
	str, err := stdio.ReadString('R')

	var row, column int
	if _, err := fmt.Sscanf(str, "\x1B[%d;%dR", &row, &column); err != nil {
		return 0, 0, err
	}

	return column, row, nil
}

func MustGetPosition() (int, int) {
	column, row, err := GetPosition()
	if err != nil {
		panic(err)
	}
	return column, row
}

func SavePosition() {
	stdio.Printf("\x1B7")
}

func RestorePosition() {
	stdio.Printf("\x1B8")
}

func Hide() {
	stdio.Printf("\x1B[?25l")
}

func Show() {
	stdio.Printf("\x1B[?25h")
}

func Save() {
	stdio.Print("\x1B7")
}

func Restore() {
	stdio.Print("\x1B8")
}
