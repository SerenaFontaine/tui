package tui

// Cell represents a single terminal cell.
type Cell struct {
	Char  rune
	Style Style
}

// emptyCell is the default cell (space with no style).
var emptyCell = Cell{Char: ' '}

// Equal returns true if two cells are visually identical.
func (c Cell) Equal(other Cell) bool {
	return c.Char == other.Char && c.Style == other.Style
}
