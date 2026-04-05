package tui

// Constraint defines how a layout slot is sized.
type Constraint struct {
	typ constraintType
	val int
}

type constraintType int

const (
	constraintFixed   constraintType = iota // Exact number of cells
	constraintPercent                       // Percentage of available space
	constraintMin                           // Minimum size, fill remaining
	constraintMax                           // Maximum size, may shrink
	constraintFlex                          // Flexible weight-based
)

// Fixed creates a constraint for an exact number of cells.
func Fixed(n int) Constraint { return Constraint{typ: constraintFixed, val: n} }

// Percent creates a constraint for a percentage of available space.
func Percent(p int) Constraint { return Constraint{typ: constraintPercent, val: p} }

// Min creates a constraint with a minimum size.
func Min(n int) Constraint { return Constraint{typ: constraintMin, val: n} }

// Max creates a constraint with a maximum size.
func Max(n int) Constraint { return Constraint{typ: constraintMax, val: n} }

// Flex creates a flexible constraint with a relative weight.
func Flex(weight int) Constraint { return Constraint{typ: constraintFlex, val: weight} }

// Direction specifies whether a layout is horizontal or vertical.
type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

// Layout splits a rectangle according to constraints.
// It returns one Rect per constraint.
func Layout(area Rect, dir Direction, constraints []Constraint) []Rect {
	if len(constraints) == 0 {
		return nil
	}

	total := area.Width
	if dir == Vertical {
		total = area.Height
	}

	sizes := resolveConstraints(constraints, total)
	rects := make([]Rect, len(sizes))

	offset := 0
	for i, size := range sizes {
		if dir == Horizontal {
			rects[i] = Rect{
				X:      area.X + offset,
				Y:      area.Y,
				Width:  size,
				Height: area.Height,
			}
		} else {
			rects[i] = Rect{
				X:      area.X,
				Y:      area.Y + offset,
				Width:  area.Width,
				Height: size,
			}
		}
		offset += size
	}

	return rects
}

// HSplit splits a rect horizontally (left to right) by constraints.
func HSplit(area Rect, constraints ...Constraint) []Rect {
	return Layout(area, Horizontal, constraints)
}

// VSplit splits a rect vertically (top to bottom) by constraints.
func VSplit(area Rect, constraints ...Constraint) []Rect {
	return Layout(area, Vertical, constraints)
}

func resolveConstraints(constraints []Constraint, total int) []int {
	n := len(constraints)
	sizes := make([]int, n)
	remaining := total
	flexTotal := 0

	// First pass: allocate fixed and percentage sizes
	for i, c := range constraints {
		switch c.typ {
		case constraintFixed:
			sizes[i] = min(c.val, remaining)
			remaining -= sizes[i]
		case constraintPercent:
			sizes[i] = min(total*c.val/100, remaining)
			remaining -= sizes[i]
		case constraintMax:
			sizes[i] = min(c.val, remaining)
			remaining -= sizes[i]
		case constraintMin:
			sizes[i] = min(c.val, remaining)
			remaining -= sizes[i]
		case constraintFlex:
			flexTotal += c.val
		}
	}

	// Second pass: distribute remaining space to flex items
	if flexTotal > 0 && remaining > 0 {
		flexRemaining := remaining
		for i, c := range constraints {
			if c.typ == constraintFlex {
				sizes[i] = flexRemaining * c.val / flexTotal
				remaining -= sizes[i]
			}
		}
		// Give any rounding remainder to the last flex item
		for i := n - 1; i >= 0; i-- {
			if constraints[i].typ == constraintFlex {
				sizes[i] += remaining
				break
			}
		}
	}

	// Ensure min constraints are satisfied from remaining flex space
	for i, c := range constraints {
		if c.typ == constraintMin && sizes[i] < c.val {
			deficit := c.val - sizes[i]
			// Try to take from flex items
			for j := n - 1; j >= 0; j-- {
				if constraints[j].typ == constraintFlex && sizes[j] > 0 {
					take := min(deficit, sizes[j])
					sizes[j] -= take
					sizes[i] += take
					deficit -= take
					if deficit == 0 {
						break
					}
				}
			}
		}
	}

	// Clamp all sizes to non-negative
	for i := range sizes {
		if sizes[i] < 0 {
			sizes[i] = 0
		}
	}

	return sizes
}
