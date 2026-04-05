package tui

import "testing"

func TestColorIsZero(t *testing.T) {
	if !NoColor.IsZero() {
		t.Error("NoColor should be zero")
	}
	if Red.IsZero() {
		t.Error("Red should not be zero")
	}
	if RGB(0, 0, 0).IsZero() {
		t.Error("RGB(0,0,0) should not be zero")
	}
}

func TestANSIColors(t *testing.T) {
	tests := []struct {
		color Color
		fg    string
	}{
		{Black, "\x1b[30m"},
		{Red, "\x1b[31m"},
		{Green, "\x1b[32m"},
		{Yellow, "\x1b[33m"},
		{Blue, "\x1b[34m"},
		{Magenta, "\x1b[35m"},
		{Cyan, "\x1b[36m"},
		{White, "\x1b[37m"},
		{BrightBlack, "\x1b[90m"},
		{BrightRed, "\x1b[91m"},
		{BrightWhite, "\x1b[97m"},
	}

	for _, tt := range tests {
		got := tt.color.fgSequence()
		if got != tt.fg {
			t.Errorf("fgSequence() = %q, want %q", got, tt.fg)
		}
	}
}

func TestANSI256Color(t *testing.T) {
	c := ANSI256(196)
	fg := c.fgSequence()
	want := "\x1b[38;5;196m"
	if fg != want {
		t.Errorf("ANSI256(196).fgSequence() = %q, want %q", fg, want)
	}

	bg := c.bgSequence()
	wantBg := "\x1b[48;5;196m"
	if bg != wantBg {
		t.Errorf("ANSI256(196).bgSequence() = %q, want %q", bg, wantBg)
	}
}

func TestRGBColor(t *testing.T) {
	c := RGB(255, 128, 0)
	fg := c.fgSequence()
	want := "\x1b[38;2;255;128;0m"
	if fg != want {
		t.Errorf("RGB(255,128,0).fgSequence() = %q, want %q", fg, want)
	}

	bg := c.bgSequence()
	wantBg := "\x1b[48;2;255;128;0m"
	if bg != wantBg {
		t.Errorf("RGB(255,128,0).bgSequence() = %q, want %q", bg, wantBg)
	}
}

func TestHexColor(t *testing.T) {
	tests := []struct {
		hex  string
		want Color
	}{
		{"#FF0000", RGB(255, 0, 0)},
		{"FF0000", RGB(255, 0, 0)},
		{"#00FF00", RGB(0, 255, 0)},
		{"invalid", NoColor},
		{"#FFF", NoColor},
	}

	for _, tt := range tests {
		got := Hex(tt.hex)
		if got != tt.want {
			t.Errorf("Hex(%q) = %v, want %v", tt.hex, got, tt.want)
		}
	}
}

func TestNoColorSequence(t *testing.T) {
	if s := NoColor.fgSequence(); s != "" {
		t.Errorf("NoColor.fgSequence() = %q, want empty", s)
	}
	if s := NoColor.bgSequence(); s != "" {
		t.Errorf("NoColor.bgSequence() = %q, want empty", s)
	}
}
