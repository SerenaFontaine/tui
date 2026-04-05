package tui

import "testing"

func TestBufferNew(t *testing.T) {
	b := NewBuffer(10, 5)
	if b.Width != 10 || b.Height != 5 {
		t.Errorf("NewBuffer(10,5) = %dx%d", b.Width, b.Height)
	}
	if len(b.Cells) != 50 {
		t.Errorf("cell count = %d, want 50", len(b.Cells))
	}
}

func TestBufferSetGet(t *testing.T) {
	b := NewBuffer(10, 10)
	b.SetChar(3, 4, 'X', NewStyle().Fg(Red))

	c := b.Get(3, 4)
	if c.Char != 'X' {
		t.Errorf("Get(3,4).Char = %q, want 'X'", c.Char)
	}
}

func TestBufferOutOfBounds(t *testing.T) {
	b := NewBuffer(5, 5)

	// Should not panic
	b.SetChar(-1, 0, 'X', Style{})
	b.SetChar(0, -1, 'X', Style{})
	b.SetChar(5, 0, 'X', Style{})
	b.SetChar(0, 5, 'X', Style{})

	c := b.Get(-1, 0)
	if c != emptyCell {
		t.Error("out-of-bounds Get should return emptyCell")
	}
}

func TestBufferSetString(t *testing.T) {
	b := NewBuffer(20, 5)
	n := b.SetString(2, 1, "hello", NewStyle())
	if n != 5 {
		t.Errorf("SetString returned %d, want 5", n)
	}

	for i, ch := range "hello" {
		c := b.Get(2+i, 1)
		if c.Char != ch {
			t.Errorf("Get(%d,1).Char = %q, want %q", 2+i, c.Char, ch)
		}
	}
}

func TestBufferSetStringTruncate(t *testing.T) {
	b := NewBuffer(5, 1)
	n := b.SetString(0, 0, "hello world", NewStyle())
	if n != 5 {
		t.Errorf("SetString should truncate to width, got %d", n)
	}
}

func TestBufferFill(t *testing.T) {
	b := NewBuffer(10, 10)
	fillCell := Cell{Char: '#', Style: NewStyle().Fg(Red)}
	b.Fill(NewRect(2, 2, 3, 3), fillCell)

	if b.Get(2, 2) != fillCell {
		t.Error("Fill should set cells in area")
	}
	if b.Get(1, 2) == fillCell {
		t.Error("Fill should not affect cells outside area")
	}
}

func TestBufferClear(t *testing.T) {
	b := NewBuffer(10, 10)
	b.SetChar(5, 5, 'X', NewStyle().Fg(Red))
	b.Clear()

	c := b.Get(5, 5)
	if c != emptyCell {
		t.Error("Clear should reset all cells to emptyCell")
	}
}

func TestBufferResize(t *testing.T) {
	b := NewBuffer(10, 10)
	b.SetChar(5, 5, 'X', Style{})
	b.Resize(20, 15)

	if b.Width != 20 || b.Height != 15 {
		t.Errorf("Resize = %dx%d, want 20x15", b.Width, b.Height)
	}
	if len(b.Cells) != 300 {
		t.Errorf("cell count after resize = %d, want 300", len(b.Cells))
	}
}

func TestBufferDiffSameContent(t *testing.T) {
	a := NewBuffer(5, 5)
	b := NewBuffer(5, 5)
	a.SetChar(2, 2, 'A', Style{})
	b.SetChar(2, 2, 'A', Style{})

	diff := b.Diff(a)
	if diff != "" {
		t.Errorf("identical buffers should produce empty diff, got %q", diff)
	}
}

func TestBufferDiffDifferentContent(t *testing.T) {
	prev := NewBuffer(5, 5)
	curr := NewBuffer(5, 5)
	curr.SetChar(2, 2, 'X', Style{})

	diff := curr.Diff(prev)
	if diff == "" {
		t.Error("different buffers should produce non-empty diff")
	}
}

func TestBufferRenderFull(t *testing.T) {
	b := NewBuffer(3, 2)
	b.SetChar(0, 0, 'A', Style{})
	b.SetChar(1, 0, 'B', Style{})

	output := b.RenderFull()
	if output == "" {
		t.Error("RenderFull should produce output")
	}
}

func TestBufferDrawHLine(t *testing.T) {
	b := NewBuffer(10, 5)
	b.DrawHLine(2, 1, 5, '-', Style{})

	for x := 2; x < 7; x++ {
		c := b.Get(x, 1)
		if c.Char != '-' {
			t.Errorf("Get(%d,1).Char = %q, want '-'", x, c.Char)
		}
	}
}

func TestBufferDrawVLine(t *testing.T) {
	b := NewBuffer(5, 10)
	b.DrawVLine(1, 2, 5, '|', Style{})

	for y := 2; y < 7; y++ {
		c := b.Get(1, y)
		if c.Char != '|' {
			t.Errorf("Get(1,%d).Char = %q, want '|'", y, c.Char)
		}
	}
}

func TestBufferMerge(t *testing.T) {
	a := NewBuffer(10, 10)
	b := NewBuffer(3, 3)
	b.SetChar(0, 0, 'X', Style{})
	b.SetChar(2, 2, 'Y', Style{})

	a.Merge(b, 5, 5)
	if a.Get(5, 5).Char != 'X' {
		t.Error("Merge should place (0,0) at offset")
	}
	if a.Get(7, 7).Char != 'Y' {
		t.Error("Merge should place (2,2) at offset")
	}
}

func TestBufferSetStringInRect(t *testing.T) {
	b := NewBuffer(20, 10)
	area := NewRect(2, 2, 5, 3)
	lines := b.SetStringInRect("hello world", area, Style{})
	if lines < 1 {
		t.Errorf("SetStringInRect should return lines used, got %d", lines)
	}
}
