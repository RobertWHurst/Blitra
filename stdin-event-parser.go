package blitra

import (
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

type EventKind int

const (
	CtrlKeyEvent EventKind = iota
	AltKeyEvent
	ShiftKeyEvent
	KeyEvent
	FocusEvent
	BlurEvent
	MouseDownEvent
	MouseUpEvent
	MouseMoveEvent
	MouseScrollEvent
	CharInputEvent
)

type Event struct {
	Kind                 EventKind
	Key                  Key
	ModifiedChar         rune
	Char                 rune
	MouseX               int
	MouseY               int
	MouseButton          MouseButton
	MouseScrollDirection MouseScrollDirection
}

type EventParser struct {
	bufLen                   int
	mx                       sync.Mutex
	buf                      []byte
	parseStallCount          int
	hasWrittenSinceLastParse bool
}

func NewEventParser() *EventParser {
	return &EventParser{
		bufLen: 1024,
		buf:    make([]byte, 0, 1024),
	}
}

func (p *EventParser) Write(buf []byte) (int, error) {
	p.mx.Lock()
	defer p.mx.Unlock()

	neededCapacity := len(p.buf) + len(buf)
	for neededCapacity > cap(p.buf) {
		p.bufLen *= 2
		newBuf := make([]byte, len(p.buf), p.bufLen)
		copy(newBuf, p.buf)
		p.buf = newBuf
	}
	prevLen := len(p.buf)
	p.buf = p.buf[:len(p.buf)+len(buf)]
	n := copy(p.buf[prevLen:], buf)
	p.hasWrittenSinceLastParse = true

	return n, nil
}

func (p *EventParser) Parse() []Event {
	p.mx.Lock()
	defer p.mx.Unlock()

	if !p.hasWrittenSinceLastParse {
		return []Event{}
	}
	p.hasWrittenSinceLastParse = false

	events := []Event{}

	// Gets the byte at the index given. Zero if out of bounds.
	b := func(i int) byte {
		if i < len(p.buf) {
			return p.buf[i]
		}
		return 0
	}

	// Checks if a byte between the given value range.
	br := func(i int, s, e byte) bool {
		if i < len(p.buf) {
			return p.buf[i] > s && p.buf[i] <= e
		}
		return false
	}

	// Checks if a string is found at the given index. Will not panic if the
	// string is longer than the buffer, but simply return false.
	sm := func(i int, s string) bool {
		if i+len(s) > len(p.buf) {
			return false
		}
		for j := 0; j < len(s); j += 1 {
			if p.buf[j+i] != s[j] {
				return false
			}
		}
		return true
	}

	// Gets an integer starting at the given index. Returns the number, and
	// the index after the number.
	n := func(i int) (int, int) {
		intStr := strings.Builder{}
		j := 0
		for ; j+i < len(p.buf) && p.buf[j+i] >= '0' && p.buf[j+i] <= '9'; j += 1 {
			intStr.WriteByte(p.buf[j+i])
		}
		if intStr.Len() == 0 {
			return 0, i
		}
		n, _ := strconv.Atoi(intStr.String())
		return n, i + j
	}

	// Takes an event pointer. If the event is not nil, it appends it to the
	// events slice. It then removes the number of bytes specified from the
	// buffer, and decrements the write offset.
	e := func(event *Event, l int) {
		if event != nil {
			events = append(events, *event)
		}
		p.buf = p.buf[l:]
		p.parseStallCount = 0
	}

	x := func() {
		if p.parseStallCount < 3 {
			p.parseStallCount += 1
			return
		}
		for i := 1; i < len(p.buf); i += 1 {
			if p.buf[i] == 0x1b {
				e(nil, i-1)
				p.parseStallCount = 0
				break
			}
		}
	}

	// Loop over stdin bytes. Collect them into events.
	for len(p.buf) != 0 {
		var event *Event

		// Escape sequence.
		if b(0) == 0x1b {

			// SS3 function key sequence.
			if b(1) == 'O' {
				switch b(2) {
				case 'P':
					event = &Event{Kind: KeyEvent, Key: F1Key}
				case 'Q':
					event = &Event{Kind: KeyEvent, Key: F2Key}
				case 'R':
					event = &Event{Kind: KeyEvent, Key: F3Key}
				case 'S':
					event = &Event{Kind: KeyEvent, Key: F4Key}
				}
				if event != nil {
					e(event, 3)
					continue
				}
			}

			// CSI sequence or SGR mouse sequence.
			if b(1) == '[' {

				// SGR mouse sequence.
				if b(2) == '<' {
					mouseButtonIndex, i := n(3)
					if b(i) == ';' {
						mouseXCoord, i := n(i + 1)
						if b(i) == ';' {
							mouseYCoord, i := n(i + 1)

							event = &Event{
								MouseButton: NoMouseButton,
								MouseX:      mouseXCoord,
								MouseY:      mouseYCoord,
							}

							switch b(i) {
							case 'M':
								event.Kind = MouseDownEvent
							case 'm':
								event.Kind = MouseUpEvent
							default:
								event = nil
							}

							if event != nil {
								switch mouseButtonIndex {
								case 0:
									event.MouseButton = LeftMouseButton
								case 1:
									event.MouseButton = MiddleMouseButton
								case 2:
									event.MouseButton = RightMouseButton
								case ' ':
									event.Kind = MouseMoveEvent
									event.MouseButton = LeftMouseButton
								case '!':
									event.Kind = MouseMoveEvent
									event.MouseButton = MiddleMouseButton
								case '"':
									event.Kind = MouseMoveEvent
									event.MouseButton = RightMouseButton
								case '#':
									event.Kind = MouseMoveEvent
								case '@':
									event.Kind = MouseScrollEvent
									event.MouseScrollDirection = MouseScrollUp
								case 'A':
									event.Kind = MouseScrollEvent
									event.MouseScrollDirection = MouseScrollDown
								default:
									event = nil
								}

								e(event, i+1)
								continue
							}
						}
					}
				}

				// CSI sequence.
				i := 2
				switch {
				case sm(i, "A"):
					event = &Event{Kind: KeyEvent, Key: UpArrowKey}
					i += 1
				case sm(i, "B"):
					event = &Event{Kind: KeyEvent, Key: DownArrowKey}
					i += 1
				case sm(i, "C"):
					event = &Event{Kind: KeyEvent, Key: RightArrowKey}
					i += 1
				case sm(i, "D"):
					event = &Event{Kind: KeyEvent, Key: LeftArrowKey}
					i += 1
				case sm(i, "F"):
					event = &Event{Kind: KeyEvent, Key: EndKey}
					i += 1
				case sm(i, "I"):
					event = &Event{Kind: FocusEvent}
					i += 1
				case sm(i, "O"):
					event = &Event{Kind: BlurEvent}
					i += 1
				case sm(i, "H"):
					event = &Event{Kind: KeyEvent, Key: HomeKey}
					i += 1
				case sm(i, "1~"):
					event = &Event{Kind: KeyEvent, Key: HomeKey}
					i += 2
				case sm(i, "2~"):
					event = &Event{Kind: KeyEvent, Key: InsertKey}
					i += 2
				case sm(i, "3~"):
					event = &Event{Kind: KeyEvent, Key: DeleteKey}
					i += 2
				case sm(i, "4~"):
					event = &Event{Kind: KeyEvent, Key: EndKey}
					i += 2
				case sm(i, "5~"):
					event = &Event{Kind: KeyEvent, Key: PageUpKey}
					i += 2
				case sm(i, "6~"):
					event = &Event{Kind: KeyEvent, Key: PageDownKey}
					i += 2
				case sm(i, "11~"):
					event = &Event{Kind: KeyEvent, Key: F1Key}
					i += 3
				case sm(i, "12~"):
					event = &Event{Kind: KeyEvent, Key: F2Key}
					i += 3
				case sm(i, "13~"):
					event = &Event{Kind: KeyEvent, Key: F3Key}
					i += 3
				case sm(i, "14~"):
					event = &Event{Kind: KeyEvent, Key: F4Key}
					i += 3
				case sm(i, "15~"):
					event = &Event{Kind: KeyEvent, Key: F5Key}
					i += 3
				case sm(i, "17~"):
					event = &Event{Kind: KeyEvent, Key: F6Key}
					i += 3
				case sm(i, "18~"):
					event = &Event{Kind: KeyEvent, Key: F7Key}
					i += 3
				case sm(i, "19~"):
					event = &Event{Kind: KeyEvent, Key: F8Key}
					i += 3
				case sm(i, "20~"):
					event = &Event{Kind: KeyEvent, Key: F9Key}
					i += 3
				case sm(i, "21~"):
					event = &Event{Kind: KeyEvent, Key: F10Key}
					i += 3
				case sm(i, "23~"):
					event = &Event{Kind: KeyEvent, Key: F11Key}
					i += 3
				case sm(i, "24~"):
					event = &Event{Kind: KeyEvent, Key: F12Key}
					i += 3

				case sm(i, "1;2A"):
					event = &Event{Kind: ShiftKeyEvent, Key: UpArrowKey}
					i += 4
				case sm(i, "1;2B"):
					event = &Event{Kind: ShiftKeyEvent, Key: DownArrowKey}
					i += 4
				case sm(i, "1;2C"):
					event = &Event{Kind: ShiftKeyEvent, Key: RightArrowKey}
					i += 4
				case sm(i, "1;2D"):
					event = &Event{Kind: ShiftKeyEvent, Key: LeftArrowKey}
					i += 4
				case sm(i, "1;2F"):
					event = &Event{Kind: ShiftKeyEvent, Key: EndKey}
					i += 4
				case sm(i, "1;2H"):
					event = &Event{Kind: ShiftKeyEvent, Key: HomeKey}
					i += 4

				case sm(i, "1;3A"):
					event = &Event{Kind: AltKeyEvent, Key: UpArrowKey}
					i += 4
				case sm(i, "1;3B"):
					event = &Event{Kind: AltKeyEvent, Key: DownArrowKey}
					i += 4
				case sm(i, "1;3C"):
					event = &Event{Kind: AltKeyEvent, Key: RightArrowKey}
					i += 4
				case sm(i, "1;3D"):
					event = &Event{Kind: AltKeyEvent, Key: LeftArrowKey}
					i += 4
				case sm(i, "1;3F"):
					event = &Event{Kind: AltKeyEvent, Key: EndKey}
					i += 4
				case sm(i, "1;3H"):
					event = &Event{Kind: AltKeyEvent, Key: HomeKey}
					i += 4

				case sm(i, "1;5A"):
					event = &Event{Kind: CtrlKeyEvent, Key: UpArrowKey}
					i += 4
				case sm(i, "1;5B"):
					event = &Event{Kind: CtrlKeyEvent, Key: DownArrowKey}
					i += 4
				case sm(i, "1;5C"):
					event = &Event{Kind: CtrlKeyEvent, Key: RightArrowKey}
					i += 4
				case sm(i, "1;5D"):
					event = &Event{Kind: CtrlKeyEvent, Key: LeftArrowKey}
					i += 4
				case sm(i, "1;5F"):
					event = &Event{Kind: CtrlKeyEvent, Key: EndKey}
					i += 4
				case sm(i, "1;5H"):
					event = &Event{Kind: CtrlKeyEvent, Key: HomeKey}
					i += 4
				}
				if event != nil {
					e(event, i)
					continue
				}

				// Unrecognized CSI sequence.
				x()
				break
			}

			// Alt/Meta key sequence.
			if b(1) != '[' && b(1) != 'O' && len(p.buf) > 1 {
				e(&Event{Kind: AltKeyEvent, ModifiedChar: rune(b(1))}, 2)
				continue
			}

			// Unrecognized escape sequence.
			x()
			break
		}

		// ASCII control characters.
		if b(0) == 0x7f || b(0) == 0x08 {
			e(&Event{Kind: KeyEvent, Key: BackspaceKey}, 1)
			continue
		}
		if b(0) == 0x0d || b(0) == 0x0a {
			e(&Event{Kind: KeyEvent, Key: EnterKey}, 1)
			continue
		}
		if br(0, 0x00, 0x1f) {
			e(&Event{Kind: CtrlKeyEvent, ModifiedChar: rune(b(0) + 0x40)}, 1)
			continue
		}

		// Skip zero bytes.
		if b(0) == 0x0 {
			e(nil, 1)
			continue
		}

		// UTF-8 compatible character input.
		r, size := utf8.DecodeRune(p.buf)
		if size != 0 {
			e(&Event{Kind: CharInputEvent, Char: r}, size)
		}
	}

	return events
}
