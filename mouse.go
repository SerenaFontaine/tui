package tui

import (
	"strconv"
	"strings"
)

// parseSGRMouse parses an SGR-mode mouse event.
// Format: CSI < Pb ; Px ; Py M/m
// Where Pb is button+modifiers, Px/Py are 1-based coordinates,
// M means press, m means release.
func parseSGRMouse(params string, final byte) (MouseMsg, bool) {
	parts := strings.Split(params, ";")
	if len(parts) != 3 {
		return MouseMsg{}, false
	}

	code, err := strconv.Atoi(parts[0])
	if err != nil {
		return MouseMsg{}, false
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return MouseMsg{}, false
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return MouseMsg{}, false
	}

	// Convert to 0-based coordinates
	x--
	y--

	msg := MouseMsg{
		X:     x,
		Y:     y,
		Alt:   code&8 != 0,
		Ctrl:  code&16 != 0,
		Shift: code&4 != 0,
	}

	if final == 'm' {
		msg.Action = MouseRelease
	} else {
		msg.Action = MousePress
	}

	if code&32 != 0 {
		msg.Action = MouseMotion
	}

	// Button
	btnCode := code & 3
	if code&64 != 0 {
		// Wheel events
		if btnCode == 0 {
			msg.Button = MouseWheelUp
		} else {
			msg.Button = MouseWheelDown
		}
		msg.Action = MousePress
	} else {
		switch btnCode {
		case 0:
			msg.Button = MouseLeft
		case 1:
			msg.Button = MouseMiddle
		case 2:
			msg.Button = MouseRight
		case 3:
			msg.Button = MouseNone
		}
	}

	return msg, true
}
