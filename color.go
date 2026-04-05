package tui

import (
	"fmt"
	"strconv"
	"strings"
)

// colorType distinguishes between color spaces.
type colorType byte

const (
	colorNone    colorType = iota
	colorANSI              // 16 standard colors (0-15)
	colorANSI256           // 256-color palette (0-255)
	colorRGB               // 24-bit true color
)

// Color represents a terminal color. The zero value is "no color" (default).
type Color struct {
	typ colorType
	val uint32 // RGB: r<<16|g<<8|b, ANSI/256: index
}

// NoColor is the zero-value color, meaning the terminal default.
var NoColor = Color{}

// Standard ANSI colors (foreground indices 0-7).
var (
	Black   = Color{typ: colorANSI, val: 0}
	Red     = Color{typ: colorANSI, val: 1}
	Green   = Color{typ: colorANSI, val: 2}
	Yellow  = Color{typ: colorANSI, val: 3}
	Blue    = Color{typ: colorANSI, val: 4}
	Magenta = Color{typ: colorANSI, val: 5}
	Cyan    = Color{typ: colorANSI, val: 6}
	White   = Color{typ: colorANSI, val: 7}
)

// Bright ANSI colors (indices 8-15).
var (
	BrightBlack   = Color{typ: colorANSI, val: 8}
	BrightRed     = Color{typ: colorANSI, val: 9}
	BrightGreen   = Color{typ: colorANSI, val: 10}
	BrightYellow  = Color{typ: colorANSI, val: 11}
	BrightBlue    = Color{typ: colorANSI, val: 12}
	BrightMagenta = Color{typ: colorANSI, val: 13}
	BrightCyan    = Color{typ: colorANSI, val: 14}
	BrightWhite   = Color{typ: colorANSI, val: 15}
)

// ANSI256 creates a color from the 256-color palette.
func ANSI256(index uint8) Color {
	return Color{typ: colorANSI256, val: uint32(index)}
}

// RGB creates a 24-bit true color.
func RGB(r, g, b uint8) Color {
	return Color{typ: colorRGB, val: uint32(r)<<16 | uint32(g)<<8 | uint32(b)}
}

// Hex creates a color from a hex string like "#FF6432" or "FF6432".
func Hex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return NoColor
	}
	val, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return NoColor
	}
	return Color{typ: colorRGB, val: uint32(val)}
}

// IsZero returns true if this is the default (no) color.
func (c Color) IsZero() bool {
	return c.typ == colorNone
}

// fgSequence returns the ANSI escape sequence for this color as foreground.
func (c Color) fgSequence() string {
	switch c.typ {
	case colorANSI:
		if c.val < 8 {
			return fmt.Sprintf("\x1b[%dm", 30+c.val)
		}
		return fmt.Sprintf("\x1b[%dm", 90+c.val-8)
	case colorANSI256:
		return fmt.Sprintf("\x1b[38;5;%dm", c.val)
	case colorRGB:
		r := (c.val >> 16) & 0xFF
		g := (c.val >> 8) & 0xFF
		b := c.val & 0xFF
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}
	return ""
}

// bgSequence returns the ANSI escape sequence for this color as background.
func (c Color) bgSequence() string {
	switch c.typ {
	case colorANSI:
		if c.val < 8 {
			return fmt.Sprintf("\x1b[%dm", 40+c.val)
		}
		return fmt.Sprintf("\x1b[%dm", 100+c.val-8)
	case colorANSI256:
		return fmt.Sprintf("\x1b[48;5;%dm", c.val)
	case colorRGB:
		r := (c.val >> 16) & 0xFF
		g := (c.val >> 8) & 0xFF
		b := c.val & 0xFF
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
	}
	return ""
}
