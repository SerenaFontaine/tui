package tui

import "fmt"

// Msg is the interface for all messages in the event system.
type Msg interface{}

// Cmd is a function that produces a message asynchronously.
// Return nil from a Cmd to produce no message.
type Cmd func() Msg

// KeyMsg represents a keyboard event.
type KeyMsg struct {
	Type KeyType
	Rune rune
	Alt  bool
}

// String returns a human-readable representation of the key event.
func (k KeyMsg) String() string {
	s := ""
	if k.Alt {
		s = "alt+"
	}
	if k.Type == KeyRune {
		return s + string(k.Rune)
	}
	return s + k.Type.String()
}

// MouseButton identifies a mouse button.
type MouseButton int

const (
	MouseLeft MouseButton = iota
	MouseMiddle
	MouseRight
	MouseWheelUp
	MouseWheelDown
	MouseNone
)

// MouseAction identifies a mouse action.
type MouseAction int

const (
	MousePress MouseAction = iota
	MouseRelease
	MouseMotion
)

// MouseMsg represents a mouse event.
type MouseMsg struct {
	X, Y   int
	Button MouseButton
	Action MouseAction
	Alt    bool
	Ctrl   bool
	Shift  bool
}

// ResizeMsg is sent when the terminal is resized.
type ResizeMsg struct {
	Width, Height int
}

// String returns a description of the resize.
func (r ResizeMsg) String() string {
	return fmt.Sprintf("resize(%d×%d)", r.Width, r.Height)
}

// QuitMsg signals the application should exit.
type QuitMsg struct{}

// FocusMsg is sent when the terminal gains focus.
type FocusMsg struct{}

// BlurMsg is sent when the terminal loses focus.
type BlurMsg struct{}

// BatchMsg wraps multiple commands to run concurrently.
type BatchMsg []Cmd

// Quit returns a command that signals the app to quit.
func Quit() Msg {
	return QuitMsg{}
}

// QuitCmd returns a Cmd that signals quit.
func QuitCmd() Cmd {
	return func() Msg { return QuitMsg{} }
}

// Batch combines multiple commands into one that runs them concurrently.
func Batch(cmds ...Cmd) Cmd {
	var valid []Cmd
	for _, c := range cmds {
		if c != nil {
			valid = append(valid, c)
		}
	}
	if len(valid) == 0 {
		return nil
	}
	if len(valid) == 1 {
		return valid[0]
	}
	return func() Msg {
		return BatchMsg(valid)
	}
}
