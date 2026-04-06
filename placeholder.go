package tui

// renderPlaceholder draws a themed placeholder box for an unsupported image.
// If alt is empty, uses "[Image]" or "[Animation]" based on the isAnimation flag.
func renderPlaceholder(buf *Buffer, area Rect, alt string, isAnimation bool, theme Theme) {
	if area.IsEmpty() {
		return
	}

	// Draw border
	block := Block{
		Border: BorderSingle,
		Style:  NewStyle().Fg(theme.BorderColor),
	}
	inner := block.Render(buf, area)

	if inner.IsEmpty() {
		return
	}

	// Fill interior with ░
	fillStyle := NewStyle().Fg(theme.TextMuted)
	fillCell := Cell{Char: '░', Style: fillStyle}
	buf.Fill(inner, fillCell)

	// Determine label
	label := alt
	if label == "" {
		if isAnimation {
			label = "Animation"
		} else {
			label = "Image"
		}
	}
	label = "[" + label + "]"

	// Truncate label if wider than inner area
	if len(label) > inner.Width {
		if inner.Width >= 3 {
			label = label[:inner.Width]
		} else {
			return // too small for any label
		}
	}

	// Center label on middle row
	midY := inner.Y + inner.Height/2
	startX := inner.X + (inner.Width-len(label))/2
	labelStyle := NewStyle().Fg(theme.TextMuted)
	buf.SetString(startX, midY, label, labelStyle)
}
