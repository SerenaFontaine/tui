package tui

import "strings"

// Style defines visual attributes for a cell.
type Style struct {
	fg, bg        Color
	bold          bool
	dim           bool
	italic        bool
	underline     bool
	blink         bool
	reverse       bool
	strikethrough bool
}

// NewStyle returns a default (empty) style.
func NewStyle() Style {
	return Style{}
}

// Fg sets the foreground color.
func (s Style) Fg(c Color) Style { s.fg = c; return s }

// Bg sets the background color.
func (s Style) Bg(c Color) Style { s.bg = c; return s }

// Bold sets bold attribute.
func (s Style) Bold(v bool) Style { s.bold = v; return s }

// Dim sets dim/faint attribute.
func (s Style) Dim(v bool) Style { s.dim = v; return s }

// Italic sets italic attribute.
func (s Style) Italic(v bool) Style { s.italic = v; return s }

// Underline sets underline attribute.
func (s Style) Underline(v bool) Style { s.underline = v; return s }

// Blink sets blink attribute.
func (s Style) Blink(v bool) Style { s.blink = v; return s }

// Reverse sets reverse video attribute.
func (s Style) Reverse(v bool) Style { s.reverse = v; return s }

// Strikethrough sets strikethrough attribute.
func (s Style) Strikethrough(v bool) Style { s.strikethrough = v; return s }

// Foreground returns the foreground color.
func (s Style) Foreground() Color { return s.fg }

// Background returns the background color.
func (s Style) Background() Color { return s.bg }

// sequence returns the full ANSI escape sequence to apply this style.
func (s Style) sequence() string {
	if s == (Style{}) {
		return ""
	}

	var b strings.Builder
	b.WriteString("\x1b[0")

	if s.bold {
		b.WriteString(";1")
	}
	if s.dim {
		b.WriteString(";2")
	}
	if s.italic {
		b.WriteString(";3")
	}
	if s.underline {
		b.WriteString(";4")
	}
	if s.blink {
		b.WriteString(";5")
	}
	if s.reverse {
		b.WriteString(";7")
	}
	if s.strikethrough {
		b.WriteString(";9")
	}

	b.WriteByte('m')

	if !s.fg.IsZero() {
		b.WriteString(s.fg.fgSequence())
	}
	if !s.bg.IsZero() {
		b.WriteString(s.bg.bgSequence())
	}

	return b.String()
}
