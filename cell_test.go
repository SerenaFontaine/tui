package tui

import "testing"

func TestCellEqual(t *testing.T) {
	a := Cell{Char: 'A', Style: NewStyle().Fg(Red)}
	b := Cell{Char: 'A', Style: NewStyle().Fg(Red)}
	if !a.Equal(b) {
		t.Error("identical cells should be equal")
	}

	c := Cell{Char: 'B', Style: NewStyle().Fg(Red)}
	if a.Equal(c) {
		t.Error("cells with different chars should not be equal")
	}

	d := Cell{Char: 'A', Style: NewStyle().Fg(Blue)}
	if a.Equal(d) {
		t.Error("cells with different styles should not be equal")
	}
}

func TestEmptyCell(t *testing.T) {
	if emptyCell.Char != ' ' {
		t.Errorf("emptyCell.Char = %q, want ' '", emptyCell.Char)
	}
	if emptyCell.Style != (Style{}) {
		t.Error("emptyCell should have zero style")
	}
}
