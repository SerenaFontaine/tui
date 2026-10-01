package tui

import "testing"

func TestSpanCreation(t *testing.T) {
	s := NewSpan("hello", NewStyle().Fg(Red))
	if s.Text != "hello" {
		t.Errorf("span text = %q, want 'hello'", s.Text)
	}
}

func TestStyledLineWidth(t *testing.T) {
	line := NewStyledLine(
		PlainSpan("hello"),
		PlainSpan(" "),
		BoldSpan("world"),
	)
	if line.Width() != 11 {
		t.Errorf("line width = %d, want 11", line.Width())
	}
}

func TestStyledLineRender(t *testing.T) {
	buf := NewBuffer(20, 5)
	line := NewStyledLine(
		ColorSpan("AB", Red),
		PlainSpan("CD"),
	)
	n := line.Render(buf, 2, 1)
	if n != 4 {
		t.Errorf("Render returned %d, want 4", n)
	}
	if buf.Get(2, 1).Char != 'A' {
		t.Errorf("(2,1) = %q, want 'A'", buf.Get(2, 1).Char)
	}
	if buf.Get(5, 1).Char != 'D' {
		t.Errorf("(5,1) = %q, want 'D'", buf.Get(5, 1).Char)
	}
}

func TestStyledLineWidthWide(t *testing.T) {
	line := NewStyledLine(PlainSpan("日本"), BoldSpan("go"))
	if line.Width() != 6 {
		t.Errorf("line width = %d, want 6", line.Width())
	}
}

func TestStyledLineRenderWide(t *testing.T) {
	buf := NewBuffer(20, 1)
	line := NewStyledLine(PlainSpan("日本"), BoldSpan("x"))
	n := line.Render(buf, 0, 0)
	if n != 5 {
		t.Errorf("Render returned %d, want 5", n)
	}

	// Each wide rune takes two cells, so the next span starts at 4
	if buf.Get(2, 0).Char != '本' {
		t.Errorf("(2,0) = %q, want '本'", buf.Get(2, 0).Char)
	}
	if buf.Get(4, 0).Char != 'x' {
		t.Errorf("(4,0) = %q, want 'x'", buf.Get(4, 0).Char)
	}
}

func TestStyledLineRenderStopsAtEdge(t *testing.T) {
	buf := NewBuffer(4, 1)
	line := NewStyledLine(PlainSpan("a日本"), PlainSpan("x"))
	n := line.Render(buf, 0, 0)

	// "a日" fills columns 0-2; "本" needs 3-4 and doesn't fit, leaving a gap
	// at 3 that "x" must not take
	if n != 3 {
		t.Errorf("Render returned %d, want 3", n)
	}
	if buf.Get(3, 0) != emptyCell {
		t.Errorf("(3,0) = %q, want emptyCell", buf.Get(3, 0).Char)
	}
}

func TestStyledTextRender(t *testing.T) {
	buf := NewBuffer(20, 5)
	text := NewStyledText(
		NewStyledLine(PlainSpan("line1")),
		NewStyledLine(PlainSpan("line2")),
	)
	area := NewRect(0, 0, 20, 5)
	text.Render(buf, area)

	if buf.Get(0, 0).Char != 'l' {
		t.Error("first line should render at y=0")
	}
	if buf.Get(0, 1).Char != 'l' {
		t.Error("second line should render at y=1")
	}
}

func TestConvenienceSpans(t *testing.T) {
	p := PlainSpan("test")
	if p.Style != (Style{}) {
		t.Error("PlainSpan should have zero style")
	}

	b := BoldSpan("test")
	if !b.Style.bold {
		t.Error("BoldSpan should be bold")
	}

	c := ColorSpan("test", Green)
	if c.Style.fg != Green {
		t.Error("ColorSpan should have specified fg color")
	}
}
