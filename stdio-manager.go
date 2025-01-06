package blitra

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/term"
)

var (
	realStdout           *os.File = os.Stdout
	fakeStdoutWriter     *os.File
	fakeStdoutReader     *os.File
	managerCount         int
	interceptedStdoutBuf bytes.Buffer
	stdioManagerGlobalMx sync.Mutex
)

// Manages the target TTY a view will be rendered to. It can be used with a
// provided TTY stdout file, or, if a TTY is not provided, it will have the
// view rendered to os.Stdout.
//
// In the event at one our more StdioManager instances are targeting
// os.Stdout, a stdout interceptor will be setup to prevent print statements
// from interfering with the rendering of the view. The captured bytes
// will be buffered.
//
// TODO:
// - Provide a way for intercepted stdout bytes to be printed to a debugging
//   component that can be rendered in the view.
// - If a debugging component is not used the intercepted stdout bytes should
//   be printed after the view is unbound.

type StdioManager struct {
	isBound            bool
	targetTTYStdout    *os.File
	prevTargetTTYState *term.State
	stdinEventParser   *EventParser
	ttySize            Size
}

// Creates a new StdioManager. If a TTY stdout file is provided, the associated
// view will be rendered to it. The target TTY will be the provided one, or
// os.Stdout if nil.
func NewStdioManager(tty *os.File) *StdioManager {
	// NOTE: The user may explicitly pass os.Stdout. If we have already
	// intercepted stdout, this reference will not be the real stdout, thus we
	// will use our own reference to the real stdout.
	if tty == nil || tty == os.Stdout {
		tty = realStdout
	}

	return &StdioManager{
		targetTTYStdout:  tty,
		stdinEventParser: NewEventParser(),
	}
}

// Indicates if the target TTY is bound and configured for rendering.
func (m *StdioManager) IsBound() bool {
	return m.isBound
}

// Binds the target TTY, setting it to use raw mode.
func (m *StdioManager) Bind() error {
	if m.isBound {
		return nil
	}
	m.isBound = true

	// Verify the target TTY of which we will be rendering the view to is indeed
	// a TTY and not a normal file or pipe.
	ttyFileDescriptor := int(m.targetTTYStdout.Fd())
	if !term.IsTerminal(ttyFileDescriptor) {
		return errors.New("cannot bind. The target is not a TTY")
	}

	// In the event that a target TTY is not provided, and a stdout interceptor
	// has not already been setup, we will setup the interceptor to prevent print
	// statements from interfering with the rendering of the view.
	if m.targetTTYStdout == realStdout && fakeStdoutReader == nil {
		stdioManagerGlobalMx.Lock()
		r, w, err := os.Pipe()
		if err != nil {
			stdioManagerGlobalMx.Unlock()
			return fmt.Errorf("failed to create pipe: %w", err)
		}

		fakeStdoutReader = r
		fakeStdoutWriter = w
		os.Stdout = fakeStdoutWriter

		managerCount += 1

		stdioManagerGlobalMx.Unlock()
	}

	// Get the screen size of the target TTY.
	if err := m.updateTTYScreenSize(); err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}

	// Set the target TTY to raw mode.
	prevTargetTTYState, err := term.MakeRaw(ttyFileDescriptor)
	if err != nil {
		return fmt.Errorf("failed to switch terminal to raw: %w", err)
	}
	m.prevTargetTTYState = prevTargetTTYState

	// Start a go routine to parse stdin events from raw mode.
	m.startStdinStdoutRoutines()

	return nil
}

// Unbinds the target TTY, setting it back to its previous state.
//
// If this is the last StdioManager instance targeting os.Stdout, the stdout
// interceptor will be removed and the captured bytes will be printed.
func (m *StdioManager) Unbind() error {
	if !m.isBound {
		return nil
	}
	m.isBound = false

	// Restore normal mode to the target TTY.
	if err := term.Restore(int(m.targetTTYStdout.Fd()), m.prevTargetTTYState); err != nil {
		return fmt.Errorf("failed to restore terminal state: %w", err)
	}
	m.prevTargetTTYState = nil

	// If the target TTY is stdout, remove the stdout interceptor.
	if m.targetTTYStdout == realStdout {
		stdioManagerGlobalMx.Lock()
		managerCount -= 1

		if managerCount == 0 {
			fakeStdoutReader.Close()
			fakeStdoutWriter.Close()
			fakeStdoutReader = nil
			fakeStdoutWriter = nil
			os.Stdout = realStdout

			// Print the captured stdout bytes.
			fmt.Print(interceptedStdoutBuf.String())
			interceptedStdoutBuf.Reset()
		}

		stdioManagerGlobalMx.Unlock()
	}

	return nil
}

func (m *StdioManager) TakeEvents() []Event {
	return m.stdinEventParser.Parse()
}

func (m *StdioManager) updateTTYScreenSize() error {
	width, height, err := term.GetSize(int(m.targetTTYStdout.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}
	m.ttySize = Size{Width: width, Height: height}
	return nil
}

func (m *StdioManager) startStdinStdoutRoutines() {
	go func() {
		stdinBuf := make([]byte, 1024)
		for m.isBound {
			m.parseStdin(stdinBuf)
		}
	}()

	go func() {
		stdoutBuf := make([]byte, 1024)
		for m.isBound {
			m.pumpFakeStdout(stdoutBuf)
		}
	}()

	go func() {
		sigWinchChan := make(chan os.Signal, 1)
		signal.Notify(sigWinchChan, syscall.SIGWINCH)
		for m.isBound {
			<-sigWinchChan
			if err := m.updateTTYScreenSize(); err != nil {
				panic("Failed to update terminal size: " + err.Error())
			}
		}
	}()
}

func (m *StdioManager) parseStdin(readBuf []byte) {
	n, err := os.Stdin.Read(readBuf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			panic("Failed to read from stdin: " + err.Error())
		}
		return
	}

	_, err = m.stdinEventParser.Write(readBuf[:n])
	if err != nil {
		panic("Failed to write to stdin event parser: " + err.Error())
	}
}

func (m *StdioManager) pumpFakeStdout(readBuf []byte) {
	n, err := fakeStdoutReader.Read(readBuf)
	if err != nil {
		if !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
			panic("Failed to read from fake stdout: " + err.Error())
		}
		return
	}

	interceptedStdoutBuf.Write(readBuf[:n])
}
