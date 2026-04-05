package widget

import "tui"

// Input is a single-line text input widget.
type Input struct {
	Value       string
	Placeholder string
	Style       tui.Style
	CursorStyle tui.Style
	Focused     bool
	Block       *tui.Block

	cursor int // cursor position in runes
	offset int // scroll offset for long text
}

// NewInput creates a new input widget.
func NewInput(placeholder string) *Input {
	return &Input{
		Placeholder: placeholder,
		CursorStyle: tui.NewStyle().Reverse(true),
	}
}

// Focus returns a Cmd that sends a FocusInputMsg.
func (in *Input) Focus() tui.Cmd {
	return func() tui.Msg {
		return FocusInputMsg{}
	}
}

// FocusInputMsg signals that the input should receive focus.
type FocusInputMsg struct{}

// SetBlock adds a border block.
func (in *Input) SetBlock(b tui.Block) *Input { in.Block = &b; return in }

// SetStyle sets the text style.
func (in *Input) SetStyle(s tui.Style) *Input { in.Style = s; return in }

// Update handles key events for the input.
func (in *Input) Update(msg tui.Msg) (*Input, tui.Cmd) {
	if !in.Focused {
		return in, nil
	}

	switch msg := msg.(type) {
	case tui.KeyMsg:
		runes := []rune(in.Value)
		switch msg.Type {
		case tui.KeyRune:
			// Insert character at cursor
			runes = insertRune(runes, in.cursor, msg.Rune)
			in.Value = string(runes)
			in.cursor++
		case tui.KeySpace:
			runes = insertRune(runes, in.cursor, ' ')
			in.Value = string(runes)
			in.cursor++
		case tui.KeyBackspace:
			if in.cursor > 0 {
				runes = append(runes[:in.cursor-1], runes[in.cursor:]...)
				in.Value = string(runes)
				in.cursor--
			}
		case tui.KeyDelete:
			if in.cursor < len(runes) {
				runes = append(runes[:in.cursor], runes[in.cursor+1:]...)
				in.Value = string(runes)
			}
		case tui.KeyLeft:
			if in.cursor > 0 {
				in.cursor--
			}
		case tui.KeyRight:
			if in.cursor < len(runes) {
				in.cursor++
			}
		case tui.KeyHome, tui.KeyCtrlA:
			in.cursor = 0
		case tui.KeyEnd, tui.KeyCtrlE:
			in.cursor = len(runes)
		case tui.KeyCtrlU:
			in.Value = string(runes[in.cursor:])
			in.cursor = 0
		case tui.KeyCtrlK:
			in.Value = string(runes[:in.cursor])
		case tui.KeyCtrlW:
			// Delete word backward
			pos := in.cursor
			for pos > 0 && runes[pos-1] == ' ' {
				pos--
			}
			for pos > 0 && runes[pos-1] != ' ' {
				pos--
			}
			in.Value = string(append(runes[:pos], runes[in.cursor:]...))
			in.cursor = pos
		}
	}

	return in, nil
}

// Render draws the input widget.
func (in *Input) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if in.Block != nil {
		inner = in.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	// Determine text to display
	displayText := in.Value
	style := in.Style

	if displayText == "" && !in.Focused {
		displayText = in.Placeholder
		style = style.Dim(true)
	}

	runes := []rune(displayText)
	visibleWidth := inner.Width

	// Adjust offset for scrolling
	if in.cursor < in.offset {
		in.offset = in.cursor
	}
	if in.cursor >= in.offset+visibleWidth {
		in.offset = in.cursor - visibleWidth + 1
	}

	// Draw visible portion
	for i := 0; i < visibleWidth && in.offset+i < len(runes); i++ {
		ch := runes[in.offset+i]
		buf.SetChar(inner.X+i, inner.Y, ch, style)
	}

	// Draw cursor
	if in.Focused {
		cursorX := inner.X + in.cursor - in.offset
		if cursorX >= inner.X && cursorX < inner.Right() {
			if in.cursor < len(runes) {
				buf.SetChar(cursorX, inner.Y, runes[in.cursor], in.CursorStyle)
			} else {
				buf.SetChar(cursorX, inner.Y, ' ', in.CursorStyle)
			}
		}
	}
}

func insertRune(runes []rune, pos int, r rune) []rune {
	if pos >= len(runes) {
		return append(runes, r)
	}
	runes = append(runes, 0)
	copy(runes[pos+1:], runes[pos:])
	runes[pos] = r
	return runes
}
