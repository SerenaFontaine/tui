package tui

import "testing"

func TestRectBasics(t *testing.T) {
	r := NewRect(5, 10, 20, 15)
	if r.Right() != 25 {
		t.Errorf("Right() = %d, want 25", r.Right())
	}
	if r.Bottom() != 25 {
		t.Errorf("Bottom() = %d, want 25", r.Bottom())
	}
	if r.Area() != 300 {
		t.Errorf("Area() = %d, want 300", r.Area())
	}
	if r.IsEmpty() {
		t.Error("non-empty rect should not be empty")
	}
}

func TestRectEmpty(t *testing.T) {
	if !NewRect(0, 0, 0, 5).IsEmpty() {
		t.Error("zero width should be empty")
	}
	if !NewRect(0, 0, 5, 0).IsEmpty() {
		t.Error("zero height should be empty")
	}
}

func TestRectContains(t *testing.T) {
	r := NewRect(10, 10, 10, 10)
	if !r.Contains(10, 10) {
		t.Error("should contain top-left corner")
	}
	if !r.Contains(19, 19) {
		t.Error("should contain bottom-right inner point")
	}
	if r.Contains(20, 20) {
		t.Error("should not contain bottom-right edge")
	}
	if r.Contains(9, 10) {
		t.Error("should not contain point left of rect")
	}
}

func TestRectInner(t *testing.T) {
	r := NewRect(0, 0, 20, 10)
	inner := r.Inner(2)
	if inner.X != 2 || inner.Y != 2 || inner.Width != 16 || inner.Height != 6 {
		t.Errorf("Inner(2) = %v, want {2,2,16,6}", inner)
	}
}

func TestRectInnerPadding(t *testing.T) {
	r := NewRect(0, 0, 20, 10)
	inner := r.InnerPadding(1, 2, 3, 4)
	if inner.X != 4 || inner.Y != 1 || inner.Width != 14 || inner.Height != 6 {
		t.Errorf("InnerPadding(1,2,3,4) = %v, want {4,1,14,6}", inner)
	}
}

func TestRectSplitVertical(t *testing.T) {
	r := NewRect(0, 0, 100, 50)
	left, right := r.SplitVertical(30)
	if left.Width != 30 {
		t.Errorf("left width = %d, want 30", left.Width)
	}
	if right.X != 30 || right.Width != 70 {
		t.Errorf("right = {X:%d W:%d}, want {X:30 W:70}", right.X, right.Width)
	}
}

func TestRectSplitHorizontal(t *testing.T) {
	r := NewRect(0, 0, 100, 50)
	top, bottom := r.SplitHorizontal(20)
	if top.Height != 20 {
		t.Errorf("top height = %d, want 20", top.Height)
	}
	if bottom.Y != 20 || bottom.Height != 30 {
		t.Errorf("bottom = {Y:%d H:%d}, want {Y:20 H:30}", bottom.Y, bottom.Height)
	}
}

func TestRectIntersect(t *testing.T) {
	a := NewRect(0, 0, 10, 10)
	b := NewRect(5, 5, 10, 10)
	inter := a.Intersect(b)
	if inter.X != 5 || inter.Y != 5 || inter.Width != 5 || inter.Height != 5 {
		t.Errorf("Intersect = %v, want {5,5,5,5}", inter)
	}
}

func TestRectIntersectNoOverlap(t *testing.T) {
	a := NewRect(0, 0, 5, 5)
	b := NewRect(10, 10, 5, 5)
	inter := a.Intersect(b)
	if !inter.IsEmpty() {
		t.Errorf("non-overlapping rects should produce empty intersection, got %v", inter)
	}
}
