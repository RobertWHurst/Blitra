package tui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/robertwhurst/blitra/draw"
	"github.com/robertwhurst/blitra/screen"
	"github.com/robertwhurst/blitra/stdio"
	"golang.org/x/crypto/ssh/terminal"
)

type ScreenBufferKind int

const (
	MainBuffer ScreenBufferKind = iota
	AltBuffer
)

type TUI struct {
	UseBuffer ScreenBufferKind
	AutoCols  bool
	AutoRows  bool
	Cols      int
	Rows      int

	usingBuffer ScreenBufferKind
	resizeChan  chan os.Signal
}

func (t *TUI) SetCursorPosition(col, row int) {
	stdio.Printf("\x1B[%d;%dH", row, column)
}

func (t *TUI) DrawAtCursor(str string) {
	t.beforeDraw()
	draw.AtCursor(str)
}

func (t *TUI) DrawAtPosition(row int, col int, str string) {
	t.beforeDraw()
	draw.AtPosition(row, col, str)
}

func (t *TUI) beforeDraw() {
	if t.usingBuffer != t.UseBuffer {
		switch t.UseBuffer {
		case MainBuffer:
			screen.MainBuffer()
		case AltBuffer:
			screen.AlternateBuffer()
			screen.ClearScrollback()
		}
		t.usingBuffer = t.UseBuffer
	}

	if t.AutoCols {
		select {
		case <-t.resizeChan:
			t.calcCols()
		default:
		}
	}
}

func (t *TUI) calcCols() {
	cols, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(fmt.Errorf("Failed to get terminal column count: %v", err))
	}
	t.Cols = cols
}

func (t *TUI) watchCols() {
	resizeChan := make(chan os.Signal, 1)
	signal.Notify(resizeChan, syscall.SIGWINCH)
	t.resizeChan = resizeChan
}
