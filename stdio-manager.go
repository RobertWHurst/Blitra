package blitra

import (
	"bytes"
	"errors"
	"fmt"
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
	ttySize            Size

	events   []Event
	eventsMx sync.Mutex
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
		targetTTYStdout: tty,
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
	m.startStdinStdoutPumps()

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
	m.eventsMx.Lock()
	events := m.events
	m.events = nil
	m.eventsMx.Unlock()
	return events
}

func (m *StdioManager) updateTTYScreenSize() error {
	width, height, err := term.GetSize(int(m.targetTTYStdout.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}
	m.ttySize = Size{Width: width, Height: height}
	return nil
}

func (m *StdioManager) startStdinStdoutPumps() {
	go func() {
		stdinBuf := make([]byte, 1024)
		for m.isBound {
			m.pumpStdin(stdinBuf)
		}
	}()
	go func() {
		stdoutBuf := make([]byte, 1024)
		for m.isBound {
			m.pumpStdout(stdoutBuf)
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

// FIXME: This does not work with more than one StdioManager instance.
// Stdin should be handled by the first StdioManager instance only, and events
// sent to all instances with their own event slices which they can clear after
// processing.
func (m *StdioManager) pumpStdin(readBuf []byte) {
	n, _ := os.Stdin.Read(readBuf)
	if n == 0 {
		return
	}

	input := readBuf[:n]

	var event Event
	switch {

	case input[0] == 0x1b && len(input) == 2:
		event = Event{Kind: CharInputEvent, Char: rune(input[1])}

	case input[0] == 0x1b:
		switch input[1] {

		case 'O':
			switch string(input[2:]) {
			case "P":
				event = Event{Kind: KeyEvent, Key: F1Key}
			case "Q":
				event = Event{Kind: KeyEvent, Key: F2Key}
			case "R":
				event = Event{Kind: KeyEvent, Key: F3Key}
			case "S":
				event = Event{Kind: KeyEvent, Key: F4Key}
			}

		case '[':
			switch string(input[2:]) {
			case "A":
				event = Event{Kind: KeyEvent, Key: UpArrowKey}
			case "B":
				event = Event{Kind: KeyEvent, Key: DownArrowKey}
			case "C":
				event = Event{Kind: KeyEvent, Key: RightArrowKey}
			case "D":
				event = Event{Kind: KeyEvent, Key: LeftArrowKey}
			case "F":
				event = Event{Kind: KeyEvent, Key: EndKey}

			case "I":
				event = Event{Kind: FocusEvent}
			case "O":
				event = Event{Kind: BlurEvent}

			case "1~", "H":
				event = Event{Kind: KeyEvent, Key: HomeKey}
			case "2~":
				event = Event{Kind: KeyEvent, Key: InsertKey}
			case "3~":
				event = Event{Kind: KeyEvent, Key: DeleteKey}
			case "4~":
				event = Event{Kind: KeyEvent, Key: EndKey}
			case "5~":
				event = Event{Kind: KeyEvent, Key: PageUpKey}
			case "6~":
				event = Event{Kind: KeyEvent, Key: PageDownKey}
			case "11~", "OP":
				event = Event{Kind: KeyEvent, Key: F1Key}
			case "12~", "OQ":
				event = Event{Kind: KeyEvent, Key: F2Key}
			case "13~", "OR":
				event = Event{Kind: KeyEvent, Key: F3Key}
			case "14~", "OS":
				event = Event{Kind: KeyEvent, Key: F4Key}
			case "15~":
				event = Event{Kind: KeyEvent, Key: F5Key}
			case "17~":
				event = Event{Kind: KeyEvent, Key: F6Key}
			case "18~":
				event = Event{Kind: KeyEvent, Key: F7Key}
			case "19~":
				event = Event{Kind: KeyEvent, Key: F8Key}
			case "20~":
				event = Event{Kind: KeyEvent, Key: F9Key}
			case "21~":
				event = Event{Kind: KeyEvent, Key: F10Key}
			case "23~":
				event = Event{Kind: KeyEvent, Key: F11Key}
			case "24~":
				event = Event{Kind: KeyEvent, Key: F12Key}

			case "1;2A":
				event = Event{Kind: ShiftKeyEvent, Key: UpArrowKey}
			case "1;2B":
				event = Event{Kind: ShiftKeyEvent, Key: DownArrowKey}
			case "1;2C":
				event = Event{Kind: ShiftKeyEvent, Key: RightArrowKey}
			case "1;2D":
				event = Event{Kind: ShiftKeyEvent, Key: LeftArrowKey}

			case "1;3A":
				event = Event{Kind: AltKeyEvent, Key: UpArrowKey}
			case "1;3B":
				event = Event{Kind: AltKeyEvent, Key: DownArrowKey}
			case "1;3C":
				event = Event{Kind: AltKeyEvent, Key: RightArrowKey}
			case "1;3D":
				event = Event{Kind: AltKeyEvent, Key: LeftArrowKey}

			case "1;5A":
				event = Event{Kind: CtrlKeyEvent, Key: UpArrowKey}
			case "1;5B":
				event = Event{Kind: CtrlKeyEvent, Key: DownArrowKey}
			case "1;5C":
				event = Event{Kind: CtrlKeyEvent, Key: RightArrowKey}
			case "1;5D":
				event = Event{Kind: CtrlKeyEvent, Key: LeftArrowKey}
			}
		}

	case input[0] == 0x7f || input[0] == 0x08:
		event = Event{Kind: KeyEvent, Key: BackspaceKey}
	case input[0] == 0x0d || input[0] == 0x0a:
		event = Event{Kind: KeyEvent, Key: EnterKey}
	case input[0] == 0x09:
		event = Event{Kind: KeyEvent, Key: TabKey}
	case input[0] == 0x20:
		event = Event{Kind: KeyEvent, Key: SpaceKey}

	case input[0] > 0x0 && input[0] <= 0x26:
		event = Event{
			Kind:         CtrlKeyEvent,
			ModifiedChar: rune(input[0] + 0x40),
		}

	default:
		event = Event{Kind: CharInputEvent, Char: rune(input[0])}
	}

	m.eventsMx.Lock()
	m.events = append(m.events, event)
	m.eventsMx.Unlock()

	// if writeToDebugFile {
	// 	debugStdinFile, err := os.OpenFile("debug-stdin.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		panic("Failed to open debug stdin file: " + err.Error())
	// 	}
	// 	switch event.Kind {
	// 	case CtrlKeyEvent:
	// 		_, err = debugStdinFile.WriteString(fmt.Sprintf("Ctrl+%c\n", event.ModifiedChar))
	// 	case AltKeyEvent:
	// 		_, err = debugStdinFile.WriteString(fmt.Sprintf("Alt+%c\n", event.ModifiedChar))
	// 	case ShiftKeyEvent:
	// 		_, err = debugStdinFile.WriteString(fmt.Sprintf("Shift+%s\n", event.Key))
	// 	case FocusEvent:
	// 		_, err = debugStdinFile.WriteString("Focus\n")
	// 	case BlurEvent:
	// 		_, err = debugStdinFile.WriteString("Blur\n")
	// 	case KeyEvent:
	// 		_, err = debugStdinFile.WriteString(fmt.Sprintf("%s\n", event.Key))
	// 	case CharInputEvent:
	// 		_, err = debugStdinFile.WriteString(fmt.Sprintf("%c\n", event.Char))
	// 	}
	// 	if err != nil {
	// 		panic("Failed to write to debug stdin file: " + err.Error())
	// 	}
	// }
}

func (m *StdioManager) pumpStdout(readBuf []byte) {
	if fakeStdoutReader == nil {
		return
	}
	n, _ := fakeStdoutReader.Read(readBuf)
	if n == 0 {
		return
	}
	interceptedStdoutBuf.Write(readBuf[:n])
}

type EventKind int

const (
	CtrlKeyEvent EventKind = iota
	AltKeyEvent
	ShiftKeyEvent
	KeyEvent
	FocusEvent
	BlurEvent
	CharInputEvent
)

type Event struct {
	Kind         EventKind
	Key          Key
	ModifiedChar rune
	Char         rune
}
