package blitra

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// Determines if there is a tty. Checks stdout by default, but can check a
// a file if provided.
func IsTTY(tty *os.File) bool {
	return term.IsTerminal(int(tty.Fd()))
}

// Gets the size of the terminal. Checks stdout by default, but can check a
// a file if provided.
func MustGetTerminalSize(tty *os.File) (int, int) {
	width, height, err := term.GetSize(int(tty.Fd()))
	if err != nil {
		panic("Failed to get terminal size: " + err.Error())
	}
	return width, height
}

func MustSwitchTTYToRaw(tty *os.File) *term.State {
	state, err := term.MakeRaw(int(tty.Fd()))
	if err != nil {
		panic("Failed to switch terminal to raw: " + err.Error())
	}
	return state
}

func MustRestoreTTYToNormal(tty *os.File, state *term.State) {
	if err := term.Restore(int(tty.Fd()), state); err != nil {
		panic("Failed to restore terminal to normal: " + err.Error())
	}
}

// TODO: Implement a stdin go routine that will collect and parse stdin data
// and provide a way to read the data from the stdio manager.

type StdioManager struct {
	realStdout *os.File
	stdout     *os.File
	stdoutSink *os.File
}

func (m *StdioManager) Bind() error {
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	m.stdout = w
	m.stdoutSink = r

	m.realStdout = os.Stdout
	os.Stdout = m.stdout

	go m.StartStdinParser()

	return nil
}

func (m *StdioManager) Unbind() {
	os.Stdout = m.realStdout
	m.stdout.Close()
	m.stdoutSink.Close()
}

func (m *StdioManager) StdoutSink() *os.File {
	return m.stdoutSink
}

func (m *StdioManager) StartStdinParser() {

}
