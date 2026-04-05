package widget

import "github.com/SerenaFontaine/tui"

// Dialog is a modal dialog box with a message and buttons.
type Dialog struct {
	Title    string
	Message  string
	Buttons  []string
	Selected int
	Style    tui.Style
	Border   tui.BorderStyle

	// Styling
	ButtonStyle         tui.Style
	SelectedButtonStyle tui.Style
	MessageStyle        tui.Style
}

// NewDialog creates a dialog with OK/Cancel buttons.
func NewDialog(title, message string) *Dialog {
	return &Dialog{
		Title:               title,
		Message:             message,
		Buttons:             []string{"OK", "Cancel"},
		Border:              tui.BorderRounded,
		ButtonStyle:         tui.NewStyle().Dim(true),
		SelectedButtonStyle: tui.NewStyle().Bold(true).Reverse(true),
	}
}

// SetButtons sets the dialog buttons.
func (d *Dialog) SetButtons(buttons ...string) *Dialog {
	d.Buttons = buttons
	return d
}

// SelectedButton returns the text of the currently selected button.
func (d *Dialog) SelectedButton() string {
	if d.Selected >= 0 && d.Selected < len(d.Buttons) {
		return d.Buttons[d.Selected]
	}
	return ""
}

// Update handles navigation between buttons.
func (d *Dialog) Update(msg tui.Msg) (*Dialog, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyLeft, tui.KeyTab:
			d.Selected--
			if d.Selected < 0 {
				d.Selected = len(d.Buttons) - 1
			}
		case tui.KeyRight, tui.KeyBacktab:
			d.Selected++
			if d.Selected >= len(d.Buttons) {
				d.Selected = 0
			}
		}
	}
	return d, nil
}

// Render draws the dialog centered in the given area.
func (d *Dialog) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	// Calculate dialog size
	msgLines := splitLines(d.Message, area.Width-8)
	dialogWidth := area.Width / 2
	if dialogWidth < 30 {
		dialogWidth = min(30, area.Width-2)
	}
	dialogHeight := len(msgLines) + 5 // borders + padding + button row

	// Center the dialog
	x := area.X + (area.Width-dialogWidth)/2
	y := area.Y + (area.Height-dialogHeight)/2
	if y < area.Y {
		y = area.Y
	}

	dialogArea := tui.NewRect(x, y, dialogWidth, dialogHeight)

	// Draw background fill
	for dy := dialogArea.Y; dy < dialogArea.Bottom(); dy++ {
		for dx := dialogArea.X; dx < dialogArea.Right(); dx++ {
			buf.SetChar(dx, dy, ' ', d.Style)
		}
	}

	// Draw border
	block := tui.Block{
		Border: d.Border,
		Title:  d.Title,
		Style:  d.Style,
	}
	inner := block.Render(buf, dialogArea)

	// Draw message
	for i, line := range msgLines {
		if i >= inner.Height-2 { // leave room for button row
			break
		}
		buf.SetString(inner.X+1, inner.Y+i, line, d.MessageStyle)
	}

	// Draw buttons at bottom
	buttonY := inner.Bottom() - 1
	totalBtnWidth := 0
	for _, btn := range d.Buttons {
		totalBtnWidth += len(btn) + 4 // [ btn ] + gap
	}
	btnX := inner.X + (inner.Width-totalBtnWidth)/2
	if btnX < inner.X {
		btnX = inner.X
	}

	for i, btn := range d.Buttons {
		label := "[ " + btn + " ]"
		style := d.ButtonStyle
		if i == d.Selected {
			style = d.SelectedButtonStyle
		}
		buf.SetString(btnX, buttonY, label, style)
		btnX += len(label) + 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
