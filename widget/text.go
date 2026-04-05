package widget

import "tui"

// Text renders static text content.
type Text struct {
	Content   string
	Style     tui.Style
	Alignment Alignment
	Block     *tui.Block
}

// Alignment controls text alignment.
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// NewText creates a new text widget.
func NewText(content string) *Text {
	return &Text{Content: content}
}

// SetStyle sets the text style.
func (t *Text) SetStyle(s tui.Style) *Text { t.Style = s; return t }

// SetAlignment sets text alignment.
func (t *Text) SetAlignment(a Alignment) *Text { t.Alignment = a; return t }

// SetBlock adds a border block around the text.
func (t *Text) SetBlock(b tui.Block) *Text { t.Block = &b; return t }

// Render draws the text into the buffer.
func (t *Text) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if t.Block != nil {
		inner = t.Block.Render(buf, area)
	}

	lines := splitLines(t.Content, inner.Width)
	for i, line := range lines {
		if i >= inner.Height {
			break
		}
		x := inner.X
		switch t.Alignment {
		case AlignCenter:
			pad := (inner.Width - runeLen(line)) / 2
			if pad > 0 {
				x += pad
			}
		case AlignRight:
			pad := inner.Width - runeLen(line)
			if pad > 0 {
				x += pad
			}
		}
		buf.SetString(x, inner.Y+i, line, t.Style)
	}
}

func splitLines(s string, maxWidth int) []string {
	if maxWidth <= 0 {
		return nil
	}
	var lines []string
	var current []rune
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, string(current))
			current = current[:0]
			continue
		}
		current = append(current, r)
		if len(current) >= maxWidth {
			lines = append(lines, string(current))
			current = current[:0]
		}
	}
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

func runeLen(s string) int {
	n := 0
	for range s {
		n++
	}
	return n
}
