package terminal

import (
	"os"

	"github.com/robertwhurst/blitra/stdio"
	"golang.org/x/term"
)

func Compatible() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func Size() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func Bell() {
	stdio.Print("\a")
}

func RawMode() (*RawModeHandle, error) {
	state, err := term.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}
	return &RawModeHandle{state: state}, nil
}

func MustRawMode() *RawModeHandle {
	h, err := RawMode()
	if err != nil {
		panic(err)
	}
	return h
}

type RawModeHandle struct {
	state *term.State
}

func (r *RawModeHandle) Restore() error {
	return term.Restore(int(os.Stdout.Fd()), r.state)
}
