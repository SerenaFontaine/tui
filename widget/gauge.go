package widget

import "github.com/SerenaFontaine/tui"

// Gauge is a full-width progress gauge with a label overlay.
type Gauge struct {
	Percent     float64
	Label       string
	Style       tui.Style
	FilledStyle tui.Style
	Block       *tui.Block
}

// NewGauge creates a new gauge widget.
func NewGauge() *Gauge {
	return &Gauge{
		FilledStyle: tui.NewStyle().Bg(tui.Blue).Fg(tui.White),
		Style:       tui.NewStyle(),
	}
}

// SetPercent sets the gauge value (0.0 to 1.0).
func (g *Gauge) SetPercent(v float64) *Gauge {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	g.Percent = v
	return g
}

// SetLabel sets the overlay label. Use "" for auto percentage.
func (g *Gauge) SetLabel(label string) *Gauge { g.Label = label; return g }

// SetBlock adds a border block.
func (g *Gauge) SetBlock(b tui.Block) *Gauge { g.Block = &b; return g }

// Render draws the gauge.
func (g *Gauge) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if g.Block != nil {
		inner = g.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	filled := int(float64(inner.Width) * g.Percent)

	// Determine label
	label := g.Label
	if label == "" {
		pct := int(g.Percent * 100)
		label = itoa(pct) + "%"
	}

	// Center the label
	labelStart := inner.X + (inner.Width-len(label))/2

	for y := inner.Y; y < inner.Bottom(); y++ {
		for x := inner.X; x < inner.Right(); x++ {
			style := g.Style
			if x-inner.X < filled {
				style = g.FilledStyle
			}

			ch := ' '
			labelIdx := x - labelStart
			if labelIdx >= 0 && labelIdx < len(label) {
				ch = rune(label[labelIdx])
			}

			buf.SetChar(x, y, ch, style)
		}
	}
}
