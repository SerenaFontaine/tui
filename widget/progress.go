package widget

import "github.com/SerenaFontaine/tui"

// Progress renders a progress bar.
type Progress struct {
	Percent     float64 // 0.0 to 1.0
	Style       tui.Style
	FilledChar  rune
	EmptyChar   rune
	FilledStyle tui.Style
	EmptyStyle  tui.Style
	Block       *tui.Block
	ShowLabel   bool
}

// NewProgress creates a new progress bar.
func NewProgress() *Progress {
	return &Progress{
		FilledChar:  '█',
		EmptyChar:   '░',
		FilledStyle: tui.NewStyle().Fg(tui.Green),
		EmptyStyle:  tui.NewStyle().Fg(tui.BrightBlack),
		ShowLabel:   true,
	}
}

// SetPercent sets the progress value (0.0 to 1.0).
func (p *Progress) SetPercent(v float64) *Progress {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	p.Percent = v
	return p
}

// SetBlock adds a border block.
func (p *Progress) SetBlock(b tui.Block) *Progress { p.Block = &b; return p }

// Render draws the progress bar.
func (p *Progress) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if p.Block != nil {
		inner = p.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	barWidth := inner.Width
	if p.ShowLabel {
		barWidth -= 5 // " 100%"
	}
	if barWidth <= 0 {
		return
	}

	filled := int(float64(barWidth) * p.Percent)
	y := inner.Y

	for x := 0; x < barWidth; x++ {
		if x < filled {
			buf.SetChar(inner.X+x, y, p.FilledChar, p.FilledStyle)
		} else {
			buf.SetChar(inner.X+x, y, p.EmptyChar, p.EmptyStyle)
		}
	}

	if p.ShowLabel {
		pct := int(p.Percent * 100)
		label := " "
		if pct < 10 {
			label += "  "
		} else if pct < 100 {
			label += " "
		}
		label += itoa(pct) + "%"
		buf.SetString(inner.X+barWidth, y, label, p.Style)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
