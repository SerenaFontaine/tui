package tui

// FocusManager tracks focus across multiple focusable components.
type FocusManager struct {
	items   []string
	current int
}

// NewFocusManager creates a focus manager with named focusable items.
func NewFocusManager(names ...string) *FocusManager {
	return &FocusManager{items: names}
}

// Current returns the name of the currently focused item.
func (f *FocusManager) Current() string {
	if f.current >= 0 && f.current < len(f.items) {
		return f.items[f.current]
	}
	return ""
}

// Index returns the index of the currently focused item.
func (f *FocusManager) Index() int {
	return f.current
}

// IsFocused returns true if the named item currently has focus.
func (f *FocusManager) IsFocused(name string) bool {
	return f.Current() == name
}

// Focus sets focus to a named item. Returns true if found.
func (f *FocusManager) Focus(name string) bool {
	for i, n := range f.items {
		if n == name {
			f.current = i
			return true
		}
	}
	return false
}

// FocusIndex sets focus by index.
func (f *FocusManager) FocusIndex(idx int) {
	if idx >= 0 && idx < len(f.items) {
		f.current = idx
	}
}

// Next moves focus to the next item (wrapping around).
func (f *FocusManager) Next() string {
	if len(f.items) == 0 {
		return ""
	}
	f.current = (f.current + 1) % len(f.items)
	return f.items[f.current]
}

// Prev moves focus to the previous item (wrapping around).
func (f *FocusManager) Prev() string {
	if len(f.items) == 0 {
		return ""
	}
	f.current--
	if f.current < 0 {
		f.current = len(f.items) - 1
	}
	return f.items[f.current]
}

// Count returns the number of focusable items.
func (f *FocusManager) Count() int {
	return len(f.items)
}

// FocusChangeMsg is sent when focus changes via the focus manager.
type FocusChangeMsg struct {
	Previous string
	Current  string
}
