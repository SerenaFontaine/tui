package tui

import (
	"strings"
	"unicode/utf8"
)

// Buffer is a 2D grid of cells representing terminal content.
type Buffer struct {
	Width  int
	Height int
	Cells  []Cell
	Images []ImagePlacement
}

// NewBuffer creates a buffer filled with empty cells.
func NewBuffer(width, height int) *Buffer {
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i] = emptyCell
	}
	return &Buffer{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// index converts (x, y) to flat index. Returns -1 if out of bounds.
func (b *Buffer) index(x, y int) int {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return -1
	}
	return y*b.Width + x
}

// Get returns the cell at (x, y). Returns emptyCell if out of bounds.
func (b *Buffer) Get(x, y int) Cell {
	if i := b.index(x, y); i >= 0 {
		return b.Cells[i]
	}
	return emptyCell
}

// Set writes a cell at (x, y). No-op if out of bounds.
func (b *Buffer) Set(x, y int, c Cell) {
	if i := b.index(x, y); i >= 0 {
		b.Cells[i] = c
	}
}

// SetChar writes a character with a style at (x, y).
func (b *Buffer) SetChar(x, y int, ch rune, style Style) {
	b.Set(x, y, Cell{Char: ch, Style: style})
}

// SetString writes a string starting at (x, y) with the given style.
// Wide runes (CJK, most emoji) take two cells.
// Returns the number of cells written.
func (b *Buffer) SetString(x, y int, s string, style Style) int {
	written := 0
	for _, ch := range s {
		w := RuneWidth(ch)
		if x+written+w > b.Width {
			break
		}
		b.setRune(x+written, y, ch, w, style)
		written += w
	}
	return written
}

// SetStringInRect writes a string within a rectangle, wrapping at boundaries.
// Returns the number of lines used.
func (b *Buffer) SetStringInRect(s string, area Rect, style Style) int {
	if area.IsEmpty() {
		return 0
	}
	x, y := area.X, area.Y
	for _, ch := range s {
		if y >= area.Bottom() {
			break
		}
		if ch == '\n' {
			x = area.X
			y++
			continue
		}
		w := RuneWidth(ch)
		if w > area.Width {
			continue // can never fit
		}
		if x+w > area.Right() {
			x = area.X
			y++
			if y >= area.Bottom() {
				break
			}
		}
		b.setRune(x, y, ch, w, style)
		x += w
	}
	return y - area.Y + 1
}

// setRune writes a rune of width w at (x, y). A wide rune also fills
// x+1 with a blank continuation cell, which Diff and RenderFull skip.
func (b *Buffer) setRune(x, y int, ch rune, w int, style Style) {
	b.SetChar(x, y, ch, style)
	if w == 2 {
		b.SetChar(x+1, y, ' ', style)
	}
}

// Fill fills a rectangular area with a cell.
func (b *Buffer) Fill(area Rect, c Cell) {
	clipped := area.Intersect(NewRect(0, 0, b.Width, b.Height))
	for y := clipped.Y; y < clipped.Bottom(); y++ {
		for x := clipped.X; x < clipped.Right(); x++ {
			b.Cells[y*b.Width+x] = c
		}
	}
}

// FillStyle fills a rectangular area with a style, preserving characters.
func (b *Buffer) FillStyle(area Rect, style Style) {
	clipped := area.Intersect(NewRect(0, 0, b.Width, b.Height))
	for y := clipped.Y; y < clipped.Bottom(); y++ {
		for x := clipped.X; x < clipped.Right(); x++ {
			i := y*b.Width + x
			b.Cells[i].Style = style
		}
	}
}

// Clear fills the entire buffer with empty cells.
func (b *Buffer) Clear() {
	for i := range b.Cells {
		b.Cells[i] = emptyCell
	}
	b.Images = b.Images[:0]
}

// AddImage registers an image placement to be rendered.
func (b *Buffer) AddImage(img ImagePlacement) {
	b.Images = append(b.Images, img)
}

// Resize changes the buffer dimensions, discarding old content.
func (b *Buffer) Resize(width, height int) {
	b.Width = width
	b.Height = height
	b.Cells = make([]Cell, width*height)
	for i := range b.Cells {
		b.Cells[i] = emptyCell
	}
	b.Images = b.Images[:0]
}

// DrawHLine draws a horizontal line from (x, y) with the given width.
func (b *Buffer) DrawHLine(x, y, width int, ch rune, style Style) {
	for i := 0; i < width; i++ {
		b.SetChar(x+i, y, ch, style)
	}
}

// DrawVLine draws a vertical line from (x, y) with the given height.
func (b *Buffer) DrawVLine(x, y, height int, ch rune, style Style) {
	for i := 0; i < height; i++ {
		b.SetChar(x, y+i, ch, style)
	}
}

// Merge copies all non-empty cells from another buffer at the given offset.
func (b *Buffer) Merge(other *Buffer, offsetX, offsetY int) {
	for y := 0; y < other.Height; y++ {
		for x := 0; x < other.Width; x++ {
			c := other.Get(x, y)
			if c != emptyCell {
				b.Set(x+offsetX, y+offsetY, c)
			}
		}
	}
	for _, img := range other.Images {
		img.X += offsetX
		img.Y += offsetY
		b.Images = append(b.Images, img)
	}
}

// Diff computes rendering commands to update the terminal from prev to current.
// Returns an ANSI escape sequence string.
func (b *Buffer) Diff(prev *Buffer) string {
	if prev == nil || prev.Width != b.Width || prev.Height != b.Height {
		return b.RenderFull()
	}

	var out strings.Builder
	var lastStyle Style
	lastX, lastY := -1, -1
	styleSet := false

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			curr := b.Get(x, y)
			old := prev.Get(x, y)
			// A wide rune covers x and x+1: redraw it if either half changed,
			// and never emit the continuation cell
			wide := x+1 < b.Width && RuneWidth(curr.Char) == 2
			if curr.Equal(old) && (!wide || b.Get(x+1, y).Equal(prev.Get(x+1, y))) {
				if wide {
					x++
				}
				continue
			}

			// Move cursor if not contiguous
			if lastX != x || lastY != y {
				out.WriteString(cursorPosition(x, y))
			}

			// Update style if changed
			if !styleSet || curr.Style != lastStyle {
				seq := curr.Style.sequence()
				if seq == "" {
					out.WriteString("\x1b[0m")
				} else {
					out.WriteString(seq)
				}
				lastStyle = curr.Style
				styleSet = true
			}

			writeRune(&out, curr.Char)
			if wide {
				x++
			}
			lastX = x + 1
			lastY = y
		}
	}

	if styleSet {
		out.WriteString("\x1b[0m")
	}

	return out.String()
}

// RenderFull outputs the entire buffer as ANSI escape sequences.
func (b *Buffer) RenderFull() string {
	var out strings.Builder
	out.Grow(b.Width * b.Height * 4) // rough estimate

	var lastStyle Style
	styleSet := false

	for y := 0; y < b.Height; y++ {
		out.WriteString(cursorPosition(0, y))
		for x := 0; x < b.Width; x++ {
			c := b.Get(x, y)
			if !styleSet || c.Style != lastStyle {
				seq := c.Style.sequence()
				if seq == "" {
					out.WriteString("\x1b[0m")
				} else {
					out.WriteString(seq)
				}
				lastStyle = c.Style
				styleSet = true
			}
			writeRune(&out, c.Char)
			if x+1 < b.Width && RuneWidth(c.Char) == 2 {
				x++ // skip the wide rune's continuation cell
			}
		}
	}

	if styleSet {
		out.WriteString("\x1b[0m")
	}

	return out.String()
}

func writeRune(b *strings.Builder, r rune) {
	if r == 0 {
		r = ' '
	}
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	b.Write(buf[:n])
}

func cursorPosition(x, y int) string {
	// ANSI cursor position is 1-based
	return "\x1b[" + itoa(y+1) + ";" + itoa(x+1) + "H"
}

func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
