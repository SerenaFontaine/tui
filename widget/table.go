package widget

import "tui"

// Table renders a data table with headers and rows.
type Table struct {
	Headers       []string
	Rows          [][]string
	ColWidths     []int // Fixed column widths (0 = auto)
	Selected      int
	Style         tui.Style
	HeaderStyle   tui.Style
	SelectedStyle tui.Style
	Block         *tui.Block
	offset        int
}

// NewTable creates a new table widget.
func NewTable(headers []string) *Table {
	return &Table{
		Headers:       headers,
		HeaderStyle:   tui.NewStyle().Bold(true).Underline(true),
		SelectedStyle: tui.NewStyle().Reverse(true),
	}
}

// SetRows sets the table data.
func (t *Table) SetRows(rows [][]string) *Table { t.Rows = rows; return t }

// SetColWidths sets fixed column widths.
func (t *Table) SetColWidths(widths []int) *Table { t.ColWidths = widths; return t }

// SetBlock adds a border block.
func (t *Table) SetBlock(b tui.Block) *Table { t.Block = &b; return t }

// SetSelectedStyle sets the selected row style.
func (t *Table) SetSelectedStyle(s tui.Style) *Table { t.SelectedStyle = s; return t }

// SelectedRow returns the selected row data, or nil if empty.
func (t *Table) SelectedRow() []string {
	if t.Selected >= 0 && t.Selected < len(t.Rows) {
		return t.Rows[t.Selected]
	}
	return nil
}

// Update handles key events for table navigation.
func (t *Table) Update(msg tui.Msg) (*Table, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyUp, tui.KeyCtrlP:
			if t.Selected > 0 {
				t.Selected--
			}
		case tui.KeyDown, tui.KeyCtrlN:
			if t.Selected < len(t.Rows)-1 {
				t.Selected++
			}
		case tui.KeyHome:
			t.Selected = 0
		case tui.KeyEnd:
			if len(t.Rows) > 0 {
				t.Selected = len(t.Rows) - 1
			}
		case tui.KeyRune:
			switch msg.Rune {
			case 'j':
				if t.Selected < len(t.Rows)-1 {
					t.Selected++
				}
			case 'k':
				if t.Selected > 0 {
					t.Selected--
				}
			}
		}
	}
	return t, nil
}

// Render draws the table.
func (t *Table) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if t.Block != nil {
		inner = t.Block.Render(buf, area)
	}

	if inner.IsEmpty() || len(t.Headers) == 0 {
		return
	}

	// Calculate column widths
	widths := t.computeWidths(inner.Width)

	// Render header
	t.renderRow(buf, inner.X, inner.Y, inner.Width, t.Headers, widths, t.HeaderStyle)

	// Adjust scroll
	contentHeight := inner.Height - 1 // minus header
	if contentHeight <= 0 {
		return
	}
	if t.Selected < t.offset {
		t.offset = t.Selected
	}
	if t.Selected >= t.offset+contentHeight {
		t.offset = t.Selected - contentHeight + 1
	}

	// Render rows
	for i := 0; i < contentHeight; i++ {
		rowIdx := t.offset + i
		if rowIdx >= len(t.Rows) {
			break
		}

		y := inner.Y + 1 + i
		style := t.Style
		if rowIdx == t.Selected {
			style = t.SelectedStyle
			for x := inner.X; x < inner.Right(); x++ {
				buf.SetChar(x, y, ' ', style)
			}
		}
		t.renderRow(buf, inner.X, y, inner.Width, t.Rows[rowIdx], widths, style)
	}
}

func (t *Table) renderRow(buf *tui.Buffer, x, y, totalWidth int, cells []string, widths []int, style tui.Style) {
	col := x
	for i, w := range widths {
		if i >= len(cells) {
			break
		}
		text := cells[i]
		if len(text) > w {
			text = text[:w]
		}
		buf.SetString(col, y, text, style)
		col += w + 1 // +1 for gap
		if col >= x+totalWidth {
			break
		}
	}
}

func (t *Table) computeWidths(totalWidth int) []int {
	n := len(t.Headers)
	widths := make([]int, n)

	if len(t.ColWidths) >= n {
		copy(widths, t.ColWidths[:n])
		return widths
	}

	// Copy any fixed widths
	fixed := 0
	autoCount := 0
	for i := 0; i < n; i++ {
		if i < len(t.ColWidths) && t.ColWidths[i] > 0 {
			widths[i] = t.ColWidths[i]
			fixed += widths[i] + 1
		} else {
			autoCount++
		}
	}

	// Distribute remaining space among auto columns
	if autoCount > 0 {
		remaining := totalWidth - fixed
		gaps := autoCount - 1
		if gaps < 0 {
			gaps = 0
		}
		remaining -= gaps
		if remaining < autoCount {
			remaining = autoCount
		}
		autoWidth := remaining / autoCount
		for i := range widths {
			if widths[i] == 0 {
				widths[i] = autoWidth
			}
		}
	}

	return widths
}
