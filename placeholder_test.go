package tui

import "testing"

func TestRenderPlaceholderDefaultLabel(t *testing.T) {
	buf := NewBuffer(20, 10)
	area := NewRect(2, 2, 16, 6)
	theme := DefaultTheme

	renderPlaceholder(buf, area, "", false, theme)

	// Border: top-left corner
	if buf.Get(2, 2).Char != BorderSingle.TopLeft {
		t.Errorf("top-left = %q, want %q", buf.Get(2, 2).Char, BorderSingle.TopLeft)
	}
	// Border: bottom-right corner
	if buf.Get(17, 7).Char != BorderSingle.BottomRight {
		t.Errorf("bottom-right = %q, want %q", buf.Get(17, 7).Char, BorderSingle.BottomRight)
	}
	// Interior fill
	if buf.Get(3, 3).Char != '░' {
		t.Errorf("interior = %q, want '░'", buf.Get(3, 3).Char)
	}
	// Centered label: "[Image]" is 7 chars, inner width is 14, so starts at 2 + 1 + (14-7)/2 = 6
	midY := 2 + 6/2 // area.Y + area.Height/2
	if buf.Get(6, midY).Char != '[' {
		t.Errorf("label start at (%d,%d) = %q, want '['", 6, midY, buf.Get(6, midY).Char)
	}
}

func TestRenderPlaceholderCustomLabel(t *testing.T) {
	buf := NewBuffer(20, 10)
	area := NewRect(0, 0, 20, 5)
	theme := DefaultTheme

	renderPlaceholder(buf, area, "Logo", false, theme)

	// Label should be "[Logo]" centered
	midY := 5 / 2 // 2
	// Inner width is 18, "[Logo]" is 6 chars, starts at 1 + (18-6)/2 = 7
	if buf.Get(7, midY).Char != '[' {
		t.Errorf("custom label start = %q, want '['", buf.Get(7, midY).Char)
	}
}

func TestRenderPlaceholderAnimation(t *testing.T) {
	buf := NewBuffer(20, 10)
	area := NewRect(0, 0, 20, 5)
	theme := DefaultTheme

	renderPlaceholder(buf, area, "", true, theme)

	// Default animation label is "[Animation]"
	midY := 5 / 2
	// Inner width is 18, "[Animation]" is 11 chars, starts at 1 + (18-11)/2 = 4
	if buf.Get(4, midY).Char != '[' {
		t.Errorf("animation label start = %q, want '['", buf.Get(4, midY).Char)
	}
}

func TestRenderPlaceholderTinyArea(t *testing.T) {
	buf := NewBuffer(5, 5)
	area := NewRect(0, 0, 3, 3)
	theme := DefaultTheme

	// Should not panic with very small area
	renderPlaceholder(buf, area, "", false, theme)

	// Border should still be drawn
	if buf.Get(0, 0).Char != BorderSingle.TopLeft {
		t.Errorf("tiny area top-left = %q, want %q", buf.Get(0, 0).Char, BorderSingle.TopLeft)
	}
}

func TestRenderPlaceholderEmptyArea(t *testing.T) {
	buf := NewBuffer(5, 5)
	area := NewRect(0, 0, 0, 0)
	theme := DefaultTheme

	// Should not panic
	renderPlaceholder(buf, area, "", false, theme)
}

func TestRenderPlaceholderMinimalArea(t *testing.T) {
	buf := NewBuffer(5, 5)
	area := NewRect(0, 0, 2, 2)
	theme := DefaultTheme

	// 2x2 has room for border only, no interior or label — should not panic
	renderPlaceholder(buf, area, "", false, theme)

	if buf.Get(0, 0).Char != BorderSingle.TopLeft {
		t.Errorf("2x2 top-left = %q, want %q", buf.Get(0, 0).Char, BorderSingle.TopLeft)
	}
}
