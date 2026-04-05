package widget

import "github.com/SerenaFontaine/tui"

// FormField represents a single field in a form.
type FormField struct {
	Label string
	Input *Input
}

// Form is a collection of labeled input fields.
type Form struct {
	Fields     []FormField
	FocusIndex int
	Style      tui.Style
	LabelStyle tui.Style
	LabelWidth int
	Block      *tui.Block
}

// NewForm creates a new form with the given fields.
func NewForm(fields ...FormField) *Form {
	// Auto-calculate label width
	maxLabel := 0
	for _, f := range fields {
		if len(f.Label) > maxLabel {
			maxLabel = len(f.Label)
		}
	}
	if len(fields) > 0 {
		fields[0].Input.Focused = true
	}
	return &Form{
		Fields:     fields,
		LabelStyle: tui.NewStyle().Bold(true),
		LabelWidth: maxLabel + 2,
	}
}

// NewFormField creates a form field with label and placeholder.
func NewFormField(label, placeholder string) FormField {
	return FormField{
		Label: label,
		Input: NewInput(placeholder),
	}
}

// SetBlock adds a border block.
func (f *Form) SetBlock(b tui.Block) *Form { f.Block = &b; return f }

// Values returns all field values as a map of label -> value.
func (f *Form) Values() map[string]string {
	vals := make(map[string]string, len(f.Fields))
	for _, field := range f.Fields {
		vals[field.Label] = field.Input.Value
	}
	return vals
}

// Value returns the value of a field by label.
func (f *Form) Value(label string) string {
	for _, field := range f.Fields {
		if field.Label == label {
			return field.Input.Value
		}
	}
	return ""
}

// FocusedField returns the currently focused field.
func (f *Form) FocusedField() *FormField {
	if f.FocusIndex >= 0 && f.FocusIndex < len(f.Fields) {
		return &f.Fields[f.FocusIndex]
	}
	return nil
}

// Update handles form navigation and input.
func (f *Form) Update(msg tui.Msg) (*Form, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyTab, tui.KeyDown:
			// Move to next field
			f.Fields[f.FocusIndex].Input.Focused = false
			f.FocusIndex = (f.FocusIndex + 1) % len(f.Fields)
			f.Fields[f.FocusIndex].Input.Focused = true
			return f, nil
		case tui.KeyBacktab, tui.KeyUp:
			// Move to previous field
			f.Fields[f.FocusIndex].Input.Focused = false
			f.FocusIndex--
			if f.FocusIndex < 0 {
				f.FocusIndex = len(f.Fields) - 1
			}
			f.Fields[f.FocusIndex].Input.Focused = true
			return f, nil
		}
	}

	// Delegate to focused input
	if f.FocusIndex >= 0 && f.FocusIndex < len(f.Fields) {
		var cmd tui.Cmd
		f.Fields[f.FocusIndex].Input, cmd = f.Fields[f.FocusIndex].Input.Update(msg)
		return f, cmd
	}

	return f, nil
}

// Render draws the form.
func (f *Form) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if f.Block != nil {
		inner = f.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	for i, field := range f.Fields {
		if i >= inner.Height {
			break
		}

		y := inner.Y + i

		// Draw label
		label := field.Label + ":"
		if len(label) > f.LabelWidth {
			label = label[:f.LabelWidth]
		}
		buf.SetString(inner.X, y, label, f.LabelStyle)

		// Draw input
		inputArea := tui.NewRect(
			inner.X+f.LabelWidth,
			y,
			inner.Width-f.LabelWidth,
			1,
		)
		field.Input.Render(buf, inputArea)
	}
}
