package tui

// Rect defines a rectangular area on screen.
type Rect struct {
	X, Y          int
	Width, Height int
}

// NewRect creates a rectangle.
func NewRect(x, y, w, h int) Rect {
	return Rect{X: x, Y: y, Width: w, Height: h}
}

// Right returns X + Width.
func (r Rect) Right() int { return r.X + r.Width }

// Bottom returns Y + Height.
func (r Rect) Bottom() int { return r.Y + r.Height }

// Area returns the total cell count.
func (r Rect) Area() int { return r.Width * r.Height }

// IsEmpty returns true if width or height is zero.
func (r Rect) IsEmpty() bool { return r.Width <= 0 || r.Height <= 0 }

// Contains returns true if the point is inside the rectangle.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.Right() && y >= r.Y && y < r.Bottom()
}

// Inner returns a rectangle with the given margin inset on all sides.
func (r Rect) Inner(margin int) Rect {
	return Rect{
		X:      r.X + margin,
		Y:      r.Y + margin,
		Width:  max(0, r.Width-2*margin),
		Height: max(0, r.Height-2*margin),
	}
}

// InnerPadding returns a rectangle inset by the given padding on each side.
func (r Rect) InnerPadding(top, right, bottom, left int) Rect {
	return Rect{
		X:      r.X + left,
		Y:      r.Y + top,
		Width:  max(0, r.Width-left-right),
		Height: max(0, r.Height-top-bottom),
	}
}

// SplitVertical splits the rect into left and right at the given x offset.
func (r Rect) SplitVertical(x int) (Rect, Rect) {
	x = clamp(x, 0, r.Width)
	left := Rect{X: r.X, Y: r.Y, Width: x, Height: r.Height}
	right := Rect{X: r.X + x, Y: r.Y, Width: r.Width - x, Height: r.Height}
	return left, right
}

// SplitHorizontal splits the rect into top and bottom at the given y offset.
func (r Rect) SplitHorizontal(y int) (Rect, Rect) {
	y = clamp(y, 0, r.Height)
	top := Rect{X: r.X, Y: r.Y, Width: r.Width, Height: y}
	bottom := Rect{X: r.X, Y: r.Y + y, Width: r.Width, Height: r.Height - y}
	return top, bottom
}

// Intersect returns the intersection of two rectangles.
func (r Rect) Intersect(other Rect) Rect {
	x := max(r.X, other.X)
	y := max(r.Y, other.Y)
	w := max(0, min(r.Right(), other.Right())-x)
	h := max(0, min(r.Bottom(), other.Bottom())-y)
	return Rect{X: x, Y: y, Width: w, Height: h}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
