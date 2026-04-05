package widget

import (
	"testing"
	"tui"
)

func newTestBuffer() *tui.Buffer {
	return tui.NewBuffer(40, 20)
}

func testArea() tui.Rect {
	return tui.NewRect(0, 0, 40, 20)
}

// --- Text ---

func TestTextRender(t *testing.T) {
	buf := newTestBuffer()
	w := NewText("Hello")
	w.Render(buf, testArea())

	if buf.Get(0, 0).Char != 'H' {
		t.Errorf("(0,0) = %q, want 'H'", buf.Get(0, 0).Char)
	}
}

func TestTextAlignment(t *testing.T) {
	buf := tui.NewBuffer(20, 5)
	area := tui.NewRect(0, 0, 20, 5)

	w := NewText("Hi").SetAlignment(AlignCenter)
	w.Render(buf, area)

	// "Hi" is 2 chars, centered in 20 → starts at (9, 0)
	if buf.Get(9, 0).Char != 'H' {
		t.Errorf("centered text: (9,0) = %q, want 'H'", buf.Get(9, 0).Char)
	}
}

func TestTextWithBlock(t *testing.T) {
	buf := tui.NewBuffer(20, 5)
	area := tui.NewRect(0, 0, 20, 5)

	block := tui.NewBlock()
	block.Title = "Test"
	w := NewText("Hi").SetBlock(block)
	w.Render(buf, area)

	// Border should be drawn
	if buf.Get(0, 0).Char != tui.BorderSingle.TopLeft {
		t.Error("block border should be rendered")
	}
	// Text inside border
	if buf.Get(1, 1).Char != 'H' {
		t.Errorf("text inside block: (1,1) = %q, want 'H'", buf.Get(1, 1).Char)
	}
}

// --- Input ---

func TestInputUpdate(t *testing.T) {
	in := NewInput("placeholder")
	in.Focused = true

	// Type "abc"
	for _, ch := range "abc" {
		in, _ = in.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: ch})
	}
	if in.Value != "abc" {
		t.Errorf("Value = %q, want 'abc'", in.Value)
	}

	// Backspace
	in, _ = in.Update(tui.KeyMsg{Type: tui.KeyBackspace})
	if in.Value != "ab" {
		t.Errorf("after backspace, Value = %q, want 'ab'", in.Value)
	}
}

func TestInputNotFocused(t *testing.T) {
	in := NewInput("")
	in.Focused = false
	in, _ = in.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'x'})
	if in.Value != "" {
		t.Error("unfocused input should not accept keys")
	}
}

func TestInputNavigation(t *testing.T) {
	in := NewInput("")
	in.Focused = true
	in.Value = "hello"
	in.cursor = 5

	in, _ = in.Update(tui.KeyMsg{Type: tui.KeyHome})
	if in.cursor != 0 {
		t.Errorf("Home: cursor = %d, want 0", in.cursor)
	}

	in, _ = in.Update(tui.KeyMsg{Type: tui.KeyEnd})
	if in.cursor != 5 {
		t.Errorf("End: cursor = %d, want 5", in.cursor)
	}
}

// --- List ---

func TestListNavigation(t *testing.T) {
	l := NewList([]string{"a", "b", "c"})
	if l.Selected != 0 {
		t.Error("initial selection should be 0")
	}

	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyDown})
	if l.Selected != 1 {
		t.Errorf("after Down, Selected = %d, want 1", l.Selected)
	}

	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyDown})
	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyDown}) // should clamp
	if l.Selected != 2 {
		t.Errorf("should clamp at 2, got %d", l.Selected)
	}

	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyUp})
	if l.Selected != 1 {
		t.Errorf("after Up, Selected = %d, want 1", l.Selected)
	}
}

func TestListSelectedItem(t *testing.T) {
	l := NewList([]string{"first", "second"})
	item := l.SelectedItem()
	if item == nil || item.Text != "first" {
		t.Error("SelectedItem should return first item")
	}
}

func TestListVimKeys(t *testing.T) {
	l := NewList([]string{"a", "b", "c"})
	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'j'})
	if l.Selected != 1 {
		t.Error("'j' should move down")
	}
	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'k'})
	if l.Selected != 0 {
		t.Error("'k' should move up")
	}
	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'G'})
	if l.Selected != 2 {
		t.Error("'G' should go to end")
	}
	l, _ = l.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'g'})
	if l.Selected != 0 {
		t.Error("'g' should go to start")
	}
}

// --- Table ---

func TestTableNavigation(t *testing.T) {
	tb := NewTable([]string{"A", "B"})
	tb.SetRows([][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}})

	tb, _ = tb.Update(tui.KeyMsg{Type: tui.KeyDown})
	if tb.Selected != 1 {
		t.Errorf("Selected = %d, want 1", tb.Selected)
	}
	row := tb.SelectedRow()
	if row[0] != "3" {
		t.Errorf("SelectedRow()[0] = %q, want '3'", row[0])
	}
}

func TestTableRender(t *testing.T) {
	buf := newTestBuffer()
	tb := NewTable([]string{"Name", "Value"})
	tb.SetRows([][]string{{"foo", "bar"}})
	tb.Render(buf, testArea())

	// Header should be rendered
	if buf.Get(0, 0).Char != 'N' {
		t.Errorf("header (0,0) = %q, want 'N'", buf.Get(0, 0).Char)
	}
}

// --- Progress ---

func TestProgressRender(t *testing.T) {
	buf := tui.NewBuffer(20, 1)
	p := NewProgress().SetPercent(0.5)
	p.Render(buf, tui.NewRect(0, 0, 20, 1))

	// Should have some filled and some empty characters
	hasFilled := false
	hasEmpty := false
	for x := 0; x < 15; x++ { // 15 = 20 - 5 (label width)
		ch := buf.Get(x, 0).Char
		if ch == '█' {
			hasFilled = true
		}
		if ch == '░' {
			hasEmpty = true
		}
	}
	if !hasFilled || !hasEmpty {
		t.Error("50% progress should have both filled and empty chars")
	}
}

// --- Tabs ---

func TestTabsNavigation(t *testing.T) {
	tabs := NewTabs([]string{"A", "B", "C"})
	if tabs.Selected != 0 {
		t.Error("initial tab should be 0")
	}

	tabs, _ = tabs.Update(tui.KeyMsg{Type: tui.KeyRight})
	if tabs.Selected != 1 {
		t.Errorf("after Right, Selected = %d, want 1", tabs.Selected)
	}

	tabs, _ = tabs.Update(tui.KeyMsg{Type: tui.KeyLeft})
	if tabs.Selected != 0 {
		t.Errorf("after Left, Selected = %d, want 0", tabs.Selected)
	}
}

func TestTabsNumberKeys(t *testing.T) {
	tabs := NewTabs([]string{"A", "B", "C"})
	tabs, _ = tabs.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: '2'})
	if tabs.Selected != 1 {
		t.Errorf("pressing '2' should select tab 1, got %d", tabs.Selected)
	}
}

// --- Spinner ---

func TestSpinnerView(t *testing.T) {
	s := NewSpinner()
	v := s.View()
	if v == "" {
		t.Error("spinner View() should not be empty")
	}
}

func TestSpinnerTick(t *testing.T) {
	s := NewSpinner()
	initial := s.View()

	// Simulate a tick
	cmd := s.Tick()
	msg := cmd()
	s, _ = s.Update(msg)

	if s.View() == initial {
		t.Error("spinner should advance after tick")
	}
}

// --- Tree ---

func TestTreeNavigation(t *testing.T) {
	tree := NewTree(
		NewTreeNode("root",
			NewTreeNode("child1"),
			NewTreeNode("child2"),
		),
	)

	if tree.SelectedNode().Text != "root" {
		t.Error("initial selection should be root")
	}

	tree, _ = tree.Update(tui.KeyMsg{Type: tui.KeyDown})
	if tree.SelectedNode().Text != "child1" {
		t.Errorf("after Down, selected = %q, want 'child1'", tree.SelectedNode().Text)
	}
}

func TestTreeCollapseExpand(t *testing.T) {
	tree := NewTree(
		NewTreeNode("root",
			NewTreeNode("child"),
		),
	)

	// Collapse root
	tree, _ = tree.Update(tui.KeyMsg{Type: tui.KeyLeft})
	// After collapse, only root should be visible
	tree, _ = tree.Update(tui.KeyMsg{Type: tui.KeyDown})
	// Should not move because there's only root
	if tree.Selected != 0 {
		t.Error("collapsed tree should have only root visible")
	}

	// Expand root
	tree, _ = tree.Update(tui.KeyMsg{Type: tui.KeyRight})
	tree, _ = tree.Update(tui.KeyMsg{Type: tui.KeyDown})
	if tree.SelectedNode().Text != "child" {
		t.Error("expanded tree should show child")
	}
}

// --- Dialog ---

func TestDialogNavigation(t *testing.T) {
	d := NewDialog("Title", "Message")
	if d.SelectedButton() != "OK" {
		t.Errorf("initial button = %q, want 'OK'", d.SelectedButton())
	}

	d, _ = d.Update(tui.KeyMsg{Type: tui.KeyRight})
	if d.SelectedButton() != "Cancel" {
		t.Errorf("after Right, button = %q, want 'Cancel'", d.SelectedButton())
	}
}

// --- Form ---

func TestFormValues(t *testing.T) {
	f := NewForm(
		NewFormField("Name", ""),
		NewFormField("Email", ""),
	)

	// Type into first field
	f, _ = f.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'J'})
	f, _ = f.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'o'})
	f, _ = f.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: 'e'})

	if f.Value("Name") != "Joe" {
		t.Errorf("Name = %q, want 'Joe'", f.Value("Name"))
	}

	// Tab to next field
	f, _ = f.Update(tui.KeyMsg{Type: tui.KeyTab})
	f, _ = f.Update(tui.KeyMsg{Type: tui.KeyRune, Rune: '@'})

	if f.Value("Email") != "@" {
		t.Errorf("Email = %q, want '@'", f.Value("Email"))
	}
}

// --- Gauge ---

func TestGaugeRender(t *testing.T) {
	buf := tui.NewBuffer(20, 1)
	g := NewGauge().SetPercent(0.5)
	g.Render(buf, tui.NewRect(0, 0, 20, 1))

	// Should have rendered something
	hasContent := false
	for x := 0; x < 20; x++ {
		if buf.Get(x, 0).Char != ' ' {
			hasContent = true
			break
		}
	}
	if !hasContent {
		t.Error("gauge should render visible content")
	}
}

// --- Viewport ---

func TestViewportScroll(t *testing.T) {
	v := NewViewport("line1\nline2\nline3\nline4\nline5")
	if v.YOffset != 0 {
		t.Error("initial offset should be 0")
	}

	v, _ = v.Update(tui.KeyMsg{Type: tui.KeyDown})
	if v.YOffset != 1 {
		t.Errorf("after Down, YOffset = %d, want 1", v.YOffset)
	}

	v, _ = v.Update(tui.KeyMsg{Type: tui.KeyUp})
	if v.YOffset != 0 {
		t.Errorf("after Up, YOffset = %d, want 0", v.YOffset)
	}
}

// --- Sparkline ---

func TestSparklineRender(t *testing.T) {
	buf := tui.NewBuffer(10, 1)
	s := NewSparkline([]float64{0.1, 0.5, 1.0, 0.3, 0.7})
	s.Render(buf, tui.NewRect(0, 0, 10, 1))

	// Should render block characters
	hasContent := false
	for x := 0; x < 5; x++ {
		if buf.Get(x, 0).Char != ' ' {
			hasContent = true
			break
		}
	}
	if !hasContent {
		t.Error("sparkline should render visible characters")
	}
}

// --- Scrollbar ---

func TestScrollbarRender(t *testing.T) {
	buf := tui.NewBuffer(1, 10)
	sb := NewScrollbar(100, 10, 0)
	sb.Render(buf, tui.NewRect(0, 0, 1, 10))

	hasThumb := false
	hasTrack := false
	for y := 0; y < 10; y++ {
		ch := buf.Get(0, y).Char
		if ch == '█' {
			hasThumb = true
		}
		if ch == '│' {
			hasTrack = true
		}
	}
	if !hasThumb || !hasTrack {
		t.Error("scrollbar should have both thumb and track characters")
	}
}

func TestScrollbarHidden(t *testing.T) {
	buf := tui.NewBuffer(1, 10)
	// Total <= visible → no scrollbar needed
	sb := NewScrollbar(5, 10, 0)
	sb.Render(buf, tui.NewRect(0, 0, 1, 10))

	// All cells should still be spaces (default)
	for y := 0; y < 10; y++ {
		if buf.Get(0, y).Char != ' ' {
			t.Error("scrollbar should not render when total <= visible")
			break
		}
	}
}
