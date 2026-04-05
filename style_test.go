package tui

import (
	"strings"
	"testing"
)

func TestStyleDefault(t *testing.T) {
	s := NewStyle()
	seq := s.sequence()
	if seq != "" {
		t.Errorf("default style sequence should be empty, got %q", seq)
	}
}

func TestStyleBold(t *testing.T) {
	s := NewStyle().Bold(true)
	seq := s.sequence()
	if !strings.Contains(seq, ";1") {
		t.Errorf("bold style should contain ';1', got %q", seq)
	}
}

func TestStyleFg(t *testing.T) {
	s := NewStyle().Fg(Red)
	seq := s.sequence()
	if !strings.Contains(seq, "\x1b[31m") {
		t.Errorf("red fg style should contain red escape, got %q", seq)
	}
}

func TestStyleBg(t *testing.T) {
	s := NewStyle().Bg(Blue)
	seq := s.sequence()
	if !strings.Contains(seq, "\x1b[44m") {
		t.Errorf("blue bg style should contain blue bg escape, got %q", seq)
	}
}

func TestStyleChaining(t *testing.T) {
	s := NewStyle().Fg(Red).Bg(Blue).Bold(true).Italic(true).Underline(true)
	seq := s.sequence()
	if !strings.Contains(seq, ";1") {
		t.Errorf("chained style missing bold")
	}
	if !strings.Contains(seq, ";3") {
		t.Errorf("chained style missing italic")
	}
	if !strings.Contains(seq, ";4") {
		t.Errorf("chained style missing underline")
	}
}

func TestStyleAccessors(t *testing.T) {
	s := NewStyle().Fg(Red).Bg(Blue)
	if s.Foreground() != Red {
		t.Error("Foreground() should return Red")
	}
	if s.Background() != Blue {
		t.Error("Background() should return Blue")
	}
}

func TestStyleAttributes(t *testing.T) {
	tests := []struct {
		name string
		s    Style
		code string
	}{
		{"dim", NewStyle().Dim(true), ";2"},
		{"blink", NewStyle().Blink(true), ";5"},
		{"reverse", NewStyle().Reverse(true), ";7"},
		{"strikethrough", NewStyle().Strikethrough(true), ";9"},
	}

	for _, tt := range tests {
		seq := tt.s.sequence()
		if !strings.Contains(seq, tt.code) {
			t.Errorf("%s: sequence %q should contain %q", tt.name, seq, tt.code)
		}
	}
}
