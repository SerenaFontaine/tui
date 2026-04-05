package widget

import "tui"

// Scrollbar renders a vertical or horizontal scrollbar.
type Scrollbar struct {
	Total   int // Total number of items/lines
	Visible int // Number visible at once
	Offset  int // Current scroll offset

	// Characters
	TrackChar  rune
	ThumbChar  rune
	Style      tui.Style
	ThumbStyle tui.Style

	Vertical bool // true = vertical (default), false = horizontal
}

// NewScrollbar creates a vertical scrollbar.
func NewScrollbar(total, visible, offset int) *Scrollbar {
	return &Scrollbar{
		Total:      total,
		Visible:    visible,
		Offset:     offset,
		TrackChar:  '│',
		ThumbChar:  '█',
		Style:      tui.NewStyle().Fg(tui.BrightBlack),
		ThumbStyle: tui.NewStyle().Fg(tui.White),
		Vertical:   true,
	}
}

// NewHScrollbar creates a horizontal scrollbar.
func NewHScrollbar(total, visible, offset int) *Scrollbar {
	s := NewScrollbar(total, visible, offset)
	s.Vertical = false
	s.TrackChar = '─'
	return s
}

// Render draws the scrollbar.
func (s *Scrollbar) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() || s.Total <= 0 || s.Total <= s.Visible {
		return
	}

	if s.Vertical {
		s.renderVertical(buf, area)
	} else {
		s.renderHorizontal(buf, area)
	}
}

func (s *Scrollbar) renderVertical(buf *tui.Buffer, area tui.Rect) {
	trackLen := area.Height
	if trackLen <= 0 {
		return
	}

	// Calculate thumb position and size
	thumbSize := max(1, trackLen*s.Visible/s.Total)
	maxOffset := s.Total - s.Visible
	if maxOffset <= 0 {
		maxOffset = 1
	}
	thumbPos := s.Offset * (trackLen - thumbSize) / maxOffset

	x := area.X
	for y := 0; y < trackLen; y++ {
		if y >= thumbPos && y < thumbPos+thumbSize {
			buf.SetChar(x, area.Y+y, s.ThumbChar, s.ThumbStyle)
		} else {
			buf.SetChar(x, area.Y+y, s.TrackChar, s.Style)
		}
	}
}

func (s *Scrollbar) renderHorizontal(buf *tui.Buffer, area tui.Rect) {
	trackLen := area.Width
	if trackLen <= 0 {
		return
	}

	thumbSize := max(1, trackLen*s.Visible/s.Total)
	maxOffset := s.Total - s.Visible
	if maxOffset <= 0 {
		maxOffset = 1
	}
	thumbPos := s.Offset * (trackLen - thumbSize) / maxOffset

	y := area.Y
	for x := 0; x < trackLen; x++ {
		if x >= thumbPos && x < thumbPos+thumbSize {
			buf.SetChar(area.X+x, y, s.ThumbChar, s.ThumbStyle)
		} else {
			buf.SetChar(area.X+x, y, s.TrackChar, s.Style)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
