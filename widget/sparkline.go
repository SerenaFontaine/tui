package widget

import "tui"

// Sparkline renders a mini line chart using braille/block characters.
type Sparkline struct {
	Data   []float64
	Style  tui.Style
	MaxVal float64 // 0 = auto
	Block  *tui.Block
}

// sparkChars contains block characters for 8 levels of height.
var sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// NewSparkline creates a new sparkline chart.
func NewSparkline(data []float64) *Sparkline {
	return &Sparkline{
		Data:  data,
		Style: tui.NewStyle().Fg(tui.Green),
	}
}

// SetData updates the sparkline data.
func (s *Sparkline) SetData(data []float64) *Sparkline {
	s.Data = data
	return s
}

// PushData appends a value, keeping the slice at maxLen.
func (s *Sparkline) PushData(val float64, maxLen int) *Sparkline {
	s.Data = append(s.Data, val)
	if len(s.Data) > maxLen {
		s.Data = s.Data[len(s.Data)-maxLen:]
	}
	return s
}

// SetMaxVal sets the maximum value for scaling. 0 = auto.
func (s *Sparkline) SetMaxVal(v float64) *Sparkline { s.MaxVal = v; return s }

// SetBlock adds a border block.
func (s *Sparkline) SetBlock(b tui.Block) *Sparkline { s.Block = &b; return s }

// SetStyle sets the sparkline style.
func (s *Sparkline) SetStyle(st tui.Style) *Sparkline { s.Style = st; return s }

// Render draws the sparkline.
func (s *Sparkline) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if s.Block != nil {
		inner = s.Block.Render(buf, area)
	}

	if inner.IsEmpty() || len(s.Data) == 0 {
		return
	}

	// Find max value
	maxVal := s.MaxVal
	if maxVal <= 0 {
		for _, v := range s.Data {
			if v > maxVal {
				maxVal = v
			}
		}
	}
	if maxVal <= 0 {
		maxVal = 1
	}

	// Determine which data points to show (rightmost)
	width := inner.Width
	startIdx := 0
	if len(s.Data) > width {
		startIdx = len(s.Data) - width
	}

	for i := 0; i < width; i++ {
		dataIdx := startIdx + i
		if dataIdx >= len(s.Data) {
			break
		}

		val := s.Data[dataIdx]
		if val < 0 {
			val = 0
		}

		// Scale to 0-7 range
		level := int(val / maxVal * float64(len(sparkChars)-1))
		if level >= len(sparkChars) {
			level = len(sparkChars) - 1
		}
		if level < 0 {
			level = 0
		}

		buf.SetChar(inner.X+i, inner.Y+inner.Height-1, sparkChars[level], s.Style)
	}
}

// SparklineGroup renders multiple sparklines stacked vertically.
type SparklineGroup struct {
	Sparklines []*Sparkline
	Block      *tui.Block
}

// NewSparklineGroup creates a group of sparklines.
func NewSparklineGroup(sparklines ...*Sparkline) *SparklineGroup {
	return &SparklineGroup{Sparklines: sparklines}
}

// SetBlock adds a border block.
func (g *SparklineGroup) SetBlock(b tui.Block) *SparklineGroup { g.Block = &b; return g }

// Render draws all sparklines stacked.
func (g *SparklineGroup) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if g.Block != nil {
		inner = g.Block.Render(buf, area)
	}

	if inner.IsEmpty() || len(g.Sparklines) == 0 {
		return
	}

	rowHeight := inner.Height / len(g.Sparklines)
	if rowHeight < 1 {
		rowHeight = 1
	}

	for i, sl := range g.Sparklines {
		y := inner.Y + i*rowHeight
		h := rowHeight
		if i == len(g.Sparklines)-1 {
			h = inner.Bottom() - y
		}
		if h <= 0 {
			break
		}
		sl.Render(buf, tui.NewRect(inner.X, y, inner.Width, h))
	}
}
