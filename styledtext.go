package tui

// Span is a run of text with a uniform style.
type Span struct {
	Text  string
	Style Style
}

// NewSpan creates a styled text span.
func NewSpan(text string, style Style) Span {
	return Span{Text: text, Style: style}
}

// StyledLine is a line composed of multiple styled spans.
type StyledLine struct {
	Spans []Span
}

// NewStyledLine creates a line from spans.
func NewStyledLine(spans ...Span) StyledLine {
	return StyledLine{Spans: spans}
}

// Width returns the total visible width of the line in cells.
func (l StyledLine) Width() int {
	w := 0
	for _, s := range l.Spans {
		for range s.Text {
			w++
		}
	}
	return w
}

// Render draws the styled line into the buffer at position (x, y).
// Returns the number of cells written.
func (l StyledLine) Render(buf *Buffer, x, y int) int {
	pos := x
	for _, span := range l.Spans {
		for _, ch := range span.Text {
			buf.SetChar(pos, y, ch, span.Style)
			pos++
		}
	}
	return pos - x
}

// StyledText is a block of styled lines for rich text display.
type StyledText struct {
	Lines []StyledLine
}

// NewStyledText creates a styled text block.
func NewStyledText(lines ...StyledLine) StyledText {
	return StyledText{Lines: lines}
}

// Render draws all lines into the buffer within the given area.
func (t StyledText) Render(buf *Buffer, area Rect) {
	for i, line := range t.Lines {
		if i >= area.Height {
			break
		}
		line.Render(buf, area.X, area.Y+i)
	}
}

// --- Convenience builders ---

// PlainSpan creates a span with default style.
func PlainSpan(text string) Span {
	return Span{Text: text}
}

// BoldSpan creates a bold span.
func BoldSpan(text string) Span {
	return Span{Text: text, Style: NewStyle().Bold(true)}
}

// ColorSpan creates a colored span.
func ColorSpan(text string, fg Color) Span {
	return Span{Text: text, Style: NewStyle().Fg(fg)}
}

// StyledSpan creates a span with the given style.
func StyledSpan(text string, s Style) Span {
	return NewSpan(text, s)
}
