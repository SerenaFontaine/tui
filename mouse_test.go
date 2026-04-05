package tui

import "testing"

func TestParseSGRMousePress(t *testing.T) {
	// Left click at (10, 20) - params are "0;11;21", final 'M'
	msg, ok := parseSGRMouse("0;11;21", 'M')
	if !ok {
		t.Fatal("parseSGRMouse should succeed")
	}
	if msg.X != 10 || msg.Y != 20 {
		t.Errorf("position = (%d,%d), want (10,20)", msg.X, msg.Y)
	}
	if msg.Button != MouseLeft {
		t.Errorf("button = %d, want MouseLeft", msg.Button)
	}
	if msg.Action != MousePress {
		t.Errorf("action = %d, want MousePress", msg.Action)
	}
}

func TestParseSGRMouseRelease(t *testing.T) {
	msg, ok := parseSGRMouse("0;5;10", 'm')
	if !ok {
		t.Fatal("should succeed")
	}
	if msg.Action != MouseRelease {
		t.Errorf("action = %d, want MouseRelease", msg.Action)
	}
}

func TestParseSGRMouseMiddle(t *testing.T) {
	msg, ok := parseSGRMouse("1;5;5", 'M')
	if !ok {
		t.Fatal("should succeed")
	}
	if msg.Button != MouseMiddle {
		t.Errorf("button = %d, want MouseMiddle", msg.Button)
	}
}

func TestParseSGRMouseRight(t *testing.T) {
	msg, ok := parseSGRMouse("2;5;5", 'M')
	if !ok {
		t.Fatal("should succeed")
	}
	if msg.Button != MouseRight {
		t.Errorf("button = %d, want MouseRight", msg.Button)
	}
}

func TestParseSGRMouseWheelUp(t *testing.T) {
	msg, ok := parseSGRMouse("64;5;5", 'M')
	if !ok {
		t.Fatal("should succeed")
	}
	if msg.Button != MouseWheelUp {
		t.Errorf("button = %d, want MouseWheelUp", msg.Button)
	}
}

func TestParseSGRMouseWheelDown(t *testing.T) {
	msg, ok := parseSGRMouse("65;5;5", 'M')
	if !ok {
		t.Fatal("should succeed")
	}
	if msg.Button != MouseWheelDown {
		t.Errorf("button = %d, want MouseWheelDown", msg.Button)
	}
}

func TestParseSGRMouseModifiers(t *testing.T) {
	// Ctrl+Alt+Left click: code = 0 | 8(alt) | 16(ctrl) = 24
	msg, ok := parseSGRMouse("24;5;5", 'M')
	if !ok {
		t.Fatal("should succeed")
	}
	if !msg.Alt {
		t.Error("Alt modifier should be set")
	}
	if !msg.Ctrl {
		t.Error("Ctrl modifier should be set")
	}
}

func TestParseSGRMouseInvalid(t *testing.T) {
	_, ok := parseSGRMouse("invalid", 'M')
	if ok {
		t.Error("invalid params should fail")
	}

	_, ok = parseSGRMouse("0;5", 'M')
	if ok {
		t.Error("too few params should fail")
	}
}
