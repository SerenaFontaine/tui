package tui

import "testing"

func TestBlockRender(t *testing.T) {
	b := NewBlock()
	buf := NewBuffer(20, 10)
	area := NewRect(0, 0, 20, 10)

	inner := b.Render(buf, area)
	if inner.X != 1 || inner.Y != 1 || inner.Width != 18 || inner.Height != 8 {
		t.Errorf("inner = %v, want {1,1,18,8}", inner)
	}

	// Check corners
	if buf.Get(0, 0).Char != BorderSingle.TopLeft {
		t.Errorf("top-left = %q, want %q", buf.Get(0, 0).Char, BorderSingle.TopLeft)
	}
	if buf.Get(19, 0).Char != BorderSingle.TopRight {
		t.Errorf("top-right = %q, want %q", buf.Get(19, 0).Char, BorderSingle.TopRight)
	}
	if buf.Get(0, 9).Char != BorderSingle.BottomLeft {
		t.Errorf("bottom-left = %q, want %q", buf.Get(0, 9).Char, BorderSingle.BottomLeft)
	}
	if buf.Get(19, 9).Char != BorderSingle.BottomRight {
		t.Errorf("bottom-right = %q, want %q", buf.Get(19, 9).Char, BorderSingle.BottomRight)
	}
}

func TestBlockRenderWithTitle(t *testing.T) {
	b := NewBlock()
	b.Title = "Test"
	buf := NewBuffer(20, 5)
	b.Render(buf, NewRect(0, 0, 20, 5))

	// Title should appear in the top border
	if buf.Get(2, 0).Char != 'T' {
		t.Errorf("title char at (2,0) = %q, want 'T'", buf.Get(2, 0).Char)
	}
	if buf.Get(3, 0).Char != 'e' {
		t.Errorf("title char at (3,0) = %q, want 'e'", buf.Get(3, 0).Char)
	}
}

func TestBlockNoBorder(t *testing.T) {
	b := Block{Border: BorderNone}
	buf := NewBuffer(10, 10)
	area := NewRect(0, 0, 10, 10)

	inner := b.Render(buf, area)
	if inner != area {
		t.Errorf("no-border inner should equal area, got %v", inner)
	}
}

func TestBlockEmptyArea(t *testing.T) {
	b := NewBlock()
	buf := NewBuffer(10, 10)
	inner := b.Render(buf, NewRect(0, 0, 0, 0))
	if !inner.IsEmpty() {
		t.Error("empty area should produce empty inner")
	}
}

func TestBorderStyles(t *testing.T) {
	styles := []struct {
		name   string
		border BorderStyle
	}{
		{"single", BorderSingle},
		{"double", BorderDouble},
		{"rounded", BorderRounded},
		{"thick", BorderThick},
		{"ascii", BorderASCII},
	}

	for _, tt := range styles {
		if tt.border.TopLeft == 0 {
			t.Errorf("%s: TopLeft should not be zero", tt.name)
		}
		if tt.border.Horizontal == 0 {
			t.Errorf("%s: Horizontal should not be zero", tt.name)
		}
		if tt.border.Vertical == 0 {
			t.Errorf("%s: Vertical should not be zero", tt.name)
		}
	}
}

func TestBlockRenderWideTitle(t *testing.T) {
	b := NewBlock()
	b.Title = "日本語"
	buf := NewBuffer(8, 3)
	b.Render(buf, NewRect(0, 0, 8, 3))

	// Room for 4 columns: "日本", then the closing space, then the corner
	if buf.Get(2, 0).Char != '日' || buf.Get(4, 0).Char != '本' {
		t.Errorf("title = %q%q at 2,4, want 日本", buf.Get(2, 0).Char, buf.Get(4, 0).Char)
	}
	if buf.Get(6, 0).Char != ' ' {
		t.Errorf("(6,0) = %q, want ' ' after title", buf.Get(6, 0).Char)
	}
	if buf.Get(7, 0).Char != BorderSingle.TopRight {
		t.Errorf("(7,0) = %q, want top-right corner", buf.Get(7, 0).Char)
	}
}

func TestBlockRenderMultibyteTitle(t *testing.T) {
	b := NewBlock()
	b.Title = "héllo"
	buf := NewBuffer(6, 3)
	b.Render(buf, NewRect(0, 0, 6, 3))

	// Room for 2 columns: "hé", with é kept whole
	if buf.Get(3, 0).Char != 'é' {
		t.Errorf("(3,0) = %q, want 'é'", buf.Get(3, 0).Char)
	}
}
