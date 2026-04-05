package tui

import "testing"

func TestKeyMsgString(t *testing.T) {
	tests := []struct {
		msg  KeyMsg
		want string
	}{
		{KeyMsg{Type: KeyRune, Rune: 'a'}, "a"},
		{KeyMsg{Type: KeyRune, Rune: 'a', Alt: true}, "alt+a"},
		{KeyMsg{Type: KeyEnter}, "enter"},
		{KeyMsg{Type: KeyCtrlC}, "ctrl+c"},
		{KeyMsg{Type: KeyUp, Alt: true}, "alt+up"},
	}

	for _, tt := range tests {
		got := tt.msg.String()
		if got != tt.want {
			t.Errorf("KeyMsg%+v.String() = %q, want %q", tt.msg, got, tt.want)
		}
	}
}

func TestResizeMsgString(t *testing.T) {
	r := ResizeMsg{Width: 80, Height: 24}
	got := r.String()
	want := "resize(80\u00d724)"
	if got != want {
		t.Errorf("ResizeMsg.String() = %q, want %q", got, want)
	}
}

func TestBatchNil(t *testing.T) {
	cmd := Batch(nil, nil)
	if cmd != nil {
		t.Error("Batch of all nils should return nil")
	}
}

func TestBatchSingle(t *testing.T) {
	called := false
	cmd := Batch(func() Msg {
		called = true
		return QuitMsg{}
	})
	if cmd == nil {
		t.Fatal("Batch of one should return that cmd")
	}
	msg := cmd()
	if !called {
		t.Error("single cmd in Batch should be called directly")
	}
	if _, ok := msg.(QuitMsg); !ok {
		t.Error("should return QuitMsg")
	}
}

func TestBatchMultiple(t *testing.T) {
	cmd := Batch(
		func() Msg { return QuitMsg{} },
		func() Msg { return QuitMsg{} },
	)
	if cmd == nil {
		t.Fatal("Batch should return a cmd")
	}
	msg := cmd()
	if _, ok := msg.(BatchMsg); !ok {
		t.Error("multi-cmd Batch should return BatchMsg")
	}
}

func TestQuit(t *testing.T) {
	msg := Quit()
	if _, ok := msg.(QuitMsg); !ok {
		t.Error("Quit() should return QuitMsg")
	}
}

func TestQuitCmd(t *testing.T) {
	cmd := QuitCmd()
	if cmd == nil {
		t.Fatal("QuitCmd() should return non-nil")
	}
	msg := cmd()
	if _, ok := msg.(QuitMsg); !ok {
		t.Error("QuitCmd()() should return QuitMsg")
	}
}

func TestKeyTypeString(t *testing.T) {
	if KeyEnter.String() != "enter" {
		t.Errorf("KeyEnter.String() = %q, want 'enter'", KeyEnter.String())
	}
	if KeyCtrlC.String() != "ctrl+c" {
		t.Errorf("KeyCtrlC.String() = %q, want 'ctrl+c'", KeyCtrlC.String())
	}
}
