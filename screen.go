package tui

import (
	"io"
	"os"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

// Screen manages the terminal: raw mode, input parsing, and output.
type Screen struct {
	in       io.Reader
	out      io.Writer
	managed  bool // true when owning the local terminal
	fd       int  // terminal fd for raw mode; -1 if unmanaged
	oldState *term.State
	sizeFunc func() (int, int)
	events   chan Msg
	done     chan struct{}
}

// newScreen creates a new screen attached to stdin/stdout.
func newScreen() *Screen {
	return &Screen{
		in:      os.Stdin,
		out:     os.Stdout,
		managed: true,
		fd:      int(os.Stdin.Fd()),
		events:  make(chan Msg, 64),
		done:    make(chan struct{}),
	}
}

// Start puts the terminal into raw mode and begins reading events.
func (s *Screen) Start() error {
	if s.managed {
		var err error
		s.oldState, err = term.MakeRaw(s.fd)
		if err != nil {
			return err
		}
	}

	// Enter alternate screen, hide cursor, enable mouse SGR tracking
	s.write("\x1b[?1049h") // alt screen
	s.write("\x1b[?25l")   // hide cursor
	s.write("\x1b[?1006h") // SGR mouse mode
	s.write("\x1b[?1003h") // all mouse tracking

	go s.readLoop()
	return nil
}

// Stop restores the terminal to its original state.
func (s *Screen) Stop() {
	close(s.done)

	s.write("\x1b[?1003l") // disable mouse tracking
	s.write("\x1b[?1006l") // disable SGR mouse
	s.write("\x1b[?25h")   // show cursor
	s.write("\x1b[?1049l") // leave alt screen

	if s.managed && s.oldState != nil {
		term.Restore(s.fd, s.oldState)
	}
}

// Suspend temporarily leaves raw mode and alternate screen for suspension.
func (s *Screen) Suspend() {
	s.write("\x1b[?1003l") // disable mouse tracking
	s.write("\x1b[?1006l") // disable SGR mouse
	s.write("\x1b[?25h")   // show cursor
	s.write("\x1b[?1049l") // leave alt screen
	if s.managed && s.oldState != nil {
		term.Restore(s.fd, s.oldState)
	}
}

// Resume re-enters raw mode and alternate screen after suspension.
func (s *Screen) Resume() {
	if s.managed {
		s.oldState, _ = term.MakeRaw(s.fd)
	}
	s.write("\x1b[?1049h") // alt screen
	s.write("\x1b[?25l")   // hide cursor
	s.write("\x1b[?1006h") // SGR mouse mode
	s.write("\x1b[?1003h") // all mouse tracking
}

// SetTitle sets the terminal window title.
func (s *Screen) SetTitle(title string) {
	s.write("\x1b]0;" + title + "\x07")
}

// ClearScreen sends a full screen clear.
func (s *Screen) ClearScreen() {
	s.write("\x1b[2J")
}

// Size returns the current terminal dimensions.
func (s *Screen) Size() (int, int) {
	if s.sizeFunc != nil {
		return s.sizeFunc()
	}
	if s.fd >= 0 {
		w, h, err := term.GetSize(s.fd)
		if err != nil {
			return 80, 24
		}
		return w, h
	}
	return 80, 24
}

// Events returns the channel of parsed events.
func (s *Screen) Events() <-chan Msg {
	return s.events
}

// Flush writes a string to the terminal output.
func (s *Screen) Flush(data string) {
	io.WriteString(s.out, data)
}

// ShowCursor shows the terminal cursor at (x, y) (0-based).
func (s *Screen) ShowCursor(x, y int) {
	s.write(cursorPosition(x, y))
	s.write("\x1b[?25h")
}

// HideCursor hides the terminal cursor.
func (s *Screen) HideCursor() {
	s.write("\x1b[?25l")
}

func (s *Screen) write(data string) {
	io.WriteString(s.out, data)
}

// readLoop continuously reads from stdin and parses into events.
func (s *Screen) readLoop() {
	buf := make([]byte, 256)
	for {
		select {
		case <-s.done:
			return
		default:
		}

		n, err := s.in.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			continue
		}

		s.parseInput(buf[:n])
	}
}

// parseInput processes raw terminal input bytes into messages.
func (s *Screen) parseInput(data []byte) {
	for len(data) > 0 {
		// ESC sequence
		if data[0] == 0x1b {
			if len(data) == 1 {
				// Lone escape - wait briefly for more data
				s.sendEvent(KeyMsg{Type: KeyEscape})
				data = data[1:]
				continue
			}

			consumed, msg := s.parseEscapeSequence(data)
			if msg != nil {
				s.sendEvent(msg)
			}
			data = data[consumed:]
			continue
		}

		// Control characters
		if data[0] < 0x20 {
			msg := s.parseControl(data[0])
			if msg != nil {
				s.sendEvent(msg)
			}
			data = data[1:]
			continue
		}

		// Space
		if data[0] == 0x20 {
			s.sendEvent(KeyMsg{Type: KeySpace})
			data = data[1:]
			continue
		}

		// DEL (backspace on some terminals)
		if data[0] == 0x7f {
			s.sendEvent(KeyMsg{Type: KeyBackspace})
			data = data[1:]
			continue
		}

		// UTF-8 character
		r, size := utf8.DecodeRune(data)
		if r != utf8.RuneError {
			s.sendEvent(KeyMsg{Type: KeyRune, Rune: r})
			data = data[size:]
		} else {
			data = data[1:]
		}
	}
}

// parseControl converts a control byte (0x00-0x1f) to a key message.
func (s *Screen) parseControl(b byte) Msg {
	switch b {
	case 0x01:
		return KeyMsg{Type: KeyCtrlA}
	case 0x02:
		return KeyMsg{Type: KeyCtrlB}
	case 0x03:
		return KeyMsg{Type: KeyCtrlC}
	case 0x04:
		return KeyMsg{Type: KeyCtrlD}
	case 0x05:
		return KeyMsg{Type: KeyCtrlE}
	case 0x06:
		return KeyMsg{Type: KeyCtrlF}
	case 0x07:
		return KeyMsg{Type: KeyCtrlG}
	case 0x08:
		return KeyMsg{Type: KeyBackspace}
	case 0x09:
		return KeyMsg{Type: KeyTab}
	case 0x0a, 0x0d:
		return KeyMsg{Type: KeyEnter}
	case 0x0b:
		return KeyMsg{Type: KeyCtrlK}
	case 0x0c:
		return KeyMsg{Type: KeyCtrlL}
	case 0x0e:
		return KeyMsg{Type: KeyCtrlN}
	case 0x0f:
		return KeyMsg{Type: KeyCtrlO}
	case 0x10:
		return KeyMsg{Type: KeyCtrlP}
	case 0x11:
		return KeyMsg{Type: KeyCtrlQ}
	case 0x12:
		return KeyMsg{Type: KeyCtrlR}
	case 0x13:
		return KeyMsg{Type: KeyCtrlS}
	case 0x14:
		return KeyMsg{Type: KeyCtrlT}
	case 0x15:
		return KeyMsg{Type: KeyCtrlU}
	case 0x16:
		return KeyMsg{Type: KeyCtrlV}
	case 0x17:
		return KeyMsg{Type: KeyCtrlW}
	case 0x18:
		return KeyMsg{Type: KeyCtrlX}
	case 0x19:
		return KeyMsg{Type: KeyCtrlY}
	case 0x1a:
		return KeyMsg{Type: KeyCtrlZ}
	}
	return nil
}

// parseEscapeSequence handles ESC-prefixed sequences.
// Returns bytes consumed and the parsed message.
func (s *Screen) parseEscapeSequence(data []byte) (int, Msg) {
	if len(data) < 2 {
		return 1, KeyMsg{Type: KeyEscape}
	}

	switch data[1] {
	case '[':
		return s.parseCSI(data)
	case 'O':
		return s.parseSS3(data)
	default:
		// Alt+key
		if data[1] >= 0x20 && data[1] < 0x7f {
			r, size := utf8.DecodeRune(data[1:])
			return 1 + size, KeyMsg{Type: KeyRune, Rune: r, Alt: true}
		}
		return 1, KeyMsg{Type: KeyEscape}
	}
}

// parseCSI handles CSI (ESC [) sequences.
func (s *Screen) parseCSI(data []byte) (int, Msg) {
	if len(data) < 3 {
		return 2, nil
	}

	// Find the end of the sequence (a byte in 0x40-0x7e range)
	i := 2
	params := ""
	for i < len(data) {
		b := data[i]
		if b >= 0x40 && b <= 0x7e {
			// Final byte
			params = string(data[2:i])
			final := b

			// SGR mouse: CSI < params M/m
			if len(data) > 2 && data[2] == '<' {
				if final == 'M' || final == 'm' {
					mouse, ok := parseSGRMouse(string(data[3:i]), final)
					if ok {
						return i + 1, mouse
					}
				}
				return i + 1, nil
			}

			return i + 1, s.csiToMsg(params, final)
		}
		i++
	}

	return len(data), nil
}

// csiToMsg converts a parsed CSI sequence to a Msg.
func (s *Screen) csiToMsg(params string, final byte) Msg {
	switch final {
	case 'A':
		return KeyMsg{Type: KeyUp}
	case 'B':
		return KeyMsg{Type: KeyDown}
	case 'C':
		return KeyMsg{Type: KeyRight}
	case 'D':
		return KeyMsg{Type: KeyLeft}
	case 'H':
		return KeyMsg{Type: KeyHome}
	case 'F':
		return KeyMsg{Type: KeyEnd}
	case 'Z':
		return KeyMsg{Type: KeyBacktab}
	case '~':
		switch params {
		case "1", "7":
			return KeyMsg{Type: KeyHome}
		case "2":
			return KeyMsg{Type: KeyInsert}
		case "3":
			return KeyMsg{Type: KeyDelete}
		case "4", "8":
			return KeyMsg{Type: KeyEnd}
		case "5":
			return KeyMsg{Type: KeyPageUp}
		case "6":
			return KeyMsg{Type: KeyPageDown}
		case "11":
			return KeyMsg{Type: KeyF1}
		case "12":
			return KeyMsg{Type: KeyF2}
		case "13":
			return KeyMsg{Type: KeyF3}
		case "14":
			return KeyMsg{Type: KeyF4}
		case "15":
			return KeyMsg{Type: KeyF5}
		case "17":
			return KeyMsg{Type: KeyF6}
		case "18":
			return KeyMsg{Type: KeyF7}
		case "19":
			return KeyMsg{Type: KeyF8}
		case "20":
			return KeyMsg{Type: KeyF9}
		case "21":
			return KeyMsg{Type: KeyF10}
		case "23":
			return KeyMsg{Type: KeyF11}
		case "24":
			return KeyMsg{Type: KeyF12}
		}
	}
	return nil
}

// parseSS3 handles SS3 (ESC O) sequences.
func (s *Screen) parseSS3(data []byte) (int, Msg) {
	if len(data) < 3 {
		return 2, nil
	}

	switch data[2] {
	case 'A':
		return 3, KeyMsg{Type: KeyUp}
	case 'B':
		return 3, KeyMsg{Type: KeyDown}
	case 'C':
		return 3, KeyMsg{Type: KeyRight}
	case 'D':
		return 3, KeyMsg{Type: KeyLeft}
	case 'H':
		return 3, KeyMsg{Type: KeyHome}
	case 'F':
		return 3, KeyMsg{Type: KeyEnd}
	case 'P':
		return 3, KeyMsg{Type: KeyF1}
	case 'Q':
		return 3, KeyMsg{Type: KeyF2}
	case 'R':
		return 3, KeyMsg{Type: KeyF3}
	case 'S':
		return 3, KeyMsg{Type: KeyF4}
	}

	return 3, nil
}

func (s *Screen) sendEvent(msg Msg) {
	select {
	case s.events <- msg:
	case <-time.After(10 * time.Millisecond):
		// Drop event if channel is full
	}
}
