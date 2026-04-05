package tui

import "testing"

func TestLayoutFixed(t *testing.T) {
	area := NewRect(0, 0, 100, 50)
	rects := VSplit(area, Fixed(10), Fixed(20), Fixed(30))

	if len(rects) != 3 {
		t.Fatalf("VSplit returned %d rects, want 3", len(rects))
	}
	if rects[0].Height != 10 {
		t.Errorf("rects[0].Height = %d, want 10", rects[0].Height)
	}
	if rects[1].Y != 10 || rects[1].Height != 20 {
		t.Errorf("rects[1] = {Y:%d H:%d}, want {Y:10 H:20}", rects[1].Y, rects[1].Height)
	}
	if rects[2].Y != 30 || rects[2].Height != 20 {
		t.Errorf("rects[2] = {Y:%d H:%d}, want {Y:30 H:20}", rects[2].Y, rects[2].Height)
	}
}

func TestLayoutFlex(t *testing.T) {
	area := NewRect(0, 0, 100, 50)
	rects := HSplit(area, Flex(1), Flex(1))

	if len(rects) != 2 {
		t.Fatalf("HSplit returned %d rects, want 2", len(rects))
	}
	if rects[0].Width != 50 {
		t.Errorf("rects[0].Width = %d, want 50", rects[0].Width)
	}
	if rects[1].X != 50 {
		t.Errorf("rects[1].X = %d, want 50", rects[1].X)
	}
}

func TestLayoutPercent(t *testing.T) {
	area := NewRect(0, 0, 200, 100)
	rects := HSplit(area, Percent(30), Flex(1))

	if len(rects) != 2 {
		t.Fatalf("HSplit returned %d rects, want 2", len(rects))
	}
	if rects[0].Width != 60 {
		t.Errorf("30%% of 200 = %d, want 60", rects[0].Width)
	}
}

func TestLayoutMixed(t *testing.T) {
	area := NewRect(0, 0, 100, 50)
	rects := VSplit(area, Fixed(3), Flex(1), Fixed(1))

	if rects[0].Height != 3 {
		t.Errorf("header height = %d, want 3", rects[0].Height)
	}
	if rects[2].Height != 1 {
		t.Errorf("footer height = %d, want 1", rects[2].Height)
	}
	if rects[1].Height != 46 {
		t.Errorf("content height = %d, want 46", rects[1].Height)
	}
}

func TestLayoutEmpty(t *testing.T) {
	area := NewRect(0, 0, 100, 50)
	rects := VSplit(area)
	if rects != nil {
		t.Error("empty constraints should return nil")
	}
}

func TestLayoutHSplit(t *testing.T) {
	area := NewRect(10, 20, 80, 40)
	rects := HSplit(area, Fixed(20), Flex(1))

	if rects[0].X != 10 {
		t.Errorf("first rect X = %d, want 10", rects[0].X)
	}
	if rects[0].Width != 20 {
		t.Errorf("first rect Width = %d, want 20", rects[0].Width)
	}
	if rects[1].X != 30 {
		t.Errorf("second rect X = %d, want 30", rects[1].X)
	}
	// All rects should have the full height
	for i, r := range rects {
		if r.Height != 40 {
			t.Errorf("rects[%d].Height = %d, want 40", i, r.Height)
		}
	}
}
