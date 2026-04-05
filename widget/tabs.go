package widget

import "tui"

// Tabs renders a tab bar with selectable tabs.
type Tabs struct {
	Titles        []string
	Selected      int
	Style         tui.Style
	ActiveStyle   tui.Style
	InactiveStyle tui.Style
	Block         *tui.Block
	Separator     string
}

// NewTabs creates a new tab bar.
func NewTabs(titles []string) *Tabs {
	return &Tabs{
		Titles:        titles,
		ActiveStyle:   tui.NewStyle().Bold(true).Underline(true),
		InactiveStyle: tui.NewStyle().Dim(true),
		Separator:     " │ ",
	}
}

// SetBlock adds a border block.
func (t *Tabs) SetBlock(b tui.Block) *Tabs { t.Block = &b; return t }

// Update handles key events for tab switching.
func (t *Tabs) Update(msg tui.Msg) (*Tabs, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyLeft:
			if t.Selected > 0 {
				t.Selected--
			}
		case tui.KeyRight:
			if t.Selected < len(t.Titles)-1 {
				t.Selected++
			}
		case tui.KeyRune:
			if msg.Rune >= '1' && msg.Rune <= '9' {
				idx := int(msg.Rune - '1')
				if idx < len(t.Titles) {
					t.Selected = idx
				}
			}
		}
	}
	return t, nil
}

// Render draws the tab bar.
func (t *Tabs) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if t.Block != nil {
		inner = t.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	x := inner.X
	for i, title := range t.Titles {
		if x >= inner.Right() {
			break
		}

		style := t.InactiveStyle
		if i == t.Selected {
			style = t.ActiveStyle
		}

		if i > 0 {
			n := buf.SetString(x, inner.Y, t.Separator, t.Style)
			x += n
		}

		n := buf.SetString(x, inner.Y, title, style)
		x += n
	}
}
