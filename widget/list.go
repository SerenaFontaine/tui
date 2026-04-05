package widget

import "tui"

// ListItem represents a single item in a list.
type ListItem struct {
	Text string
	Data any
}

// List is a scrollable, selectable list widget.
type List struct {
	Items         []ListItem
	Selected      int
	Style         tui.Style
	SelectedStyle tui.Style
	Block         *tui.Block
	offset        int // scroll offset
}

// NewList creates a list from string items.
func NewList(items []string) *List {
	listItems := make([]ListItem, len(items))
	for i, s := range items {
		listItems[i] = ListItem{Text: s}
	}
	return &List{
		Items:         listItems,
		SelectedStyle: tui.NewStyle().Reverse(true),
	}
}

// NewListFromItems creates a list from ListItem values.
func NewListFromItems(items []ListItem) *List {
	return &List{
		Items:         items,
		SelectedStyle: tui.NewStyle().Reverse(true),
	}
}

// SetBlock adds a border block.
func (l *List) SetBlock(b tui.Block) *List { l.Block = &b; return l }

// SetStyle sets the default item style.
func (l *List) SetStyle(s tui.Style) *List { l.Style = s; return l }

// SetSelectedStyle sets the selected item style.
func (l *List) SetSelectedStyle(s tui.Style) *List { l.SelectedStyle = s; return l }

// SelectedItem returns the currently selected item, or nil if the list is empty.
func (l *List) SelectedItem() *ListItem {
	if l.Selected >= 0 && l.Selected < len(l.Items) {
		return &l.Items[l.Selected]
	}
	return nil
}

// Update handles key events for navigation.
func (l *List) Update(msg tui.Msg) (*List, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyUp, tui.KeyCtrlP:
			if l.Selected > 0 {
				l.Selected--
			}
		case tui.KeyDown, tui.KeyCtrlN:
			if l.Selected < len(l.Items)-1 {
				l.Selected++
			}
		case tui.KeyHome:
			l.Selected = 0
		case tui.KeyEnd:
			l.Selected = len(l.Items) - 1
		case tui.KeyPageUp:
			l.Selected -= 10
			if l.Selected < 0 {
				l.Selected = 0
			}
		case tui.KeyPageDown:
			l.Selected += 10
			if l.Selected >= len(l.Items) {
				l.Selected = len(l.Items) - 1
			}
		case tui.KeyRune:
			switch msg.Rune {
			case 'j':
				if l.Selected < len(l.Items)-1 {
					l.Selected++
				}
			case 'k':
				if l.Selected > 0 {
					l.Selected--
				}
			case 'g':
				l.Selected = 0
			case 'G':
				l.Selected = len(l.Items) - 1
			}
		}
	}

	return l, nil
}

// Render draws the list.
func (l *List) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if l.Block != nil {
		inner = l.Block.Render(buf, area)
	}

	if inner.IsEmpty() || len(l.Items) == 0 {
		return
	}

	// Adjust scroll offset to keep selected item visible
	visibleHeight := inner.Height
	if l.Selected < l.offset {
		l.offset = l.Selected
	}
	if l.Selected >= l.offset+visibleHeight {
		l.offset = l.Selected - visibleHeight + 1
	}

	// Render visible items
	for i := 0; i < visibleHeight; i++ {
		idx := l.offset + i
		if idx >= len(l.Items) {
			break
		}

		item := l.Items[idx]
		style := l.Style
		if idx == l.Selected {
			style = l.SelectedStyle
			// Fill the entire line with selected style
			for x := inner.X; x < inner.Right(); x++ {
				buf.SetChar(x, inner.Y+i, ' ', style)
			}
		}

		text := item.Text
		if len(text) > inner.Width {
			text = text[:inner.Width]
		}
		buf.SetString(inner.X, inner.Y+i, text, style)
	}
}
