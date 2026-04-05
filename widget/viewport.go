package widget

import (
	"strings"
	"github.com/SerenaFontaine/tui"
)

// Viewport is a scrollable text display widget.
type Viewport struct {
	Content string
	Style   tui.Style
	Block   *tui.Block
	YOffset int

	lines []string // cached split lines
	width int      // width at which lines were split
}

// NewViewport creates a new viewport.
func NewViewport(content string) *Viewport {
	return &Viewport{Content: content}
}

// SetContent updates the displayed text.
func (v *Viewport) SetContent(content string) *Viewport {
	v.Content = content
	v.lines = nil // invalidate cache
	return v
}

// SetBlock adds a border block.
func (v *Viewport) SetBlock(b tui.Block) *Viewport { v.Block = &b; return v }

// SetStyle sets the text style.
func (v *Viewport) SetStyle(s tui.Style) *Viewport { v.Style = s; return v }

// ScrollTo scrolls to an absolute line offset.
func (v *Viewport) ScrollTo(y int) {
	if y < 0 {
		y = 0
	}
	v.YOffset = y
}

// LineCount returns the total number of wrapped lines.
func (v *Viewport) LineCount() int {
	return len(v.lines)
}

// Update handles scrolling keys.
func (v *Viewport) Update(msg tui.Msg) (*Viewport, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyUp, tui.KeyCtrlP:
			if v.YOffset > 0 {
				v.YOffset--
			}
		case tui.KeyDown, tui.KeyCtrlN:
			v.YOffset++
		case tui.KeyPageUp:
			v.YOffset -= 20
			if v.YOffset < 0 {
				v.YOffset = 0
			}
		case tui.KeyPageDown:
			v.YOffset += 20
		case tui.KeyHome:
			v.YOffset = 0
		case tui.KeyEnd:
			v.YOffset = len(v.lines) - 1
			if v.YOffset < 0 {
				v.YOffset = 0
			}
		case tui.KeyRune:
			switch msg.Rune {
			case 'j':
				v.YOffset++
			case 'k':
				if v.YOffset > 0 {
					v.YOffset--
				}
			}
		}
	case tui.MouseMsg:
		switch msg.Button {
		case tui.MouseWheelUp:
			if v.YOffset > 0 {
				v.YOffset -= 3
				if v.YOffset < 0 {
					v.YOffset = 0
				}
			}
		case tui.MouseWheelDown:
			v.YOffset += 3
		}
	}

	return v, nil
}

// Render draws the viewport content.
func (v *Viewport) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if v.Block != nil {
		inner = v.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	// Rebuild line cache if needed
	if v.lines == nil || v.width != inner.Width {
		v.width = inner.Width
		v.lines = wrapText(v.Content, inner.Width)
	}

	// Clamp offset
	maxOffset := len(v.lines) - inner.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if v.YOffset > maxOffset {
		v.YOffset = maxOffset
	}

	// Render visible lines
	for i := 0; i < inner.Height; i++ {
		lineIdx := v.YOffset + i
		if lineIdx >= len(v.lines) {
			break
		}
		line := v.lines[lineIdx]
		buf.SetString(inner.X, inner.Y+i, line, v.Style)
	}
}

func wrapText(s string, width int) []string {
	if width <= 0 {
		return nil
	}

	rawLines := strings.Split(s, "\n")
	var lines []string

	for _, raw := range rawLines {
		if len(raw) == 0 {
			lines = append(lines, "")
			continue
		}
		runes := []rune(raw)
		for len(runes) > 0 {
			end := width
			if end > len(runes) {
				end = len(runes)
			}
			lines = append(lines, string(runes[:end]))
			runes = runes[end:]
		}
	}

	return lines
}
