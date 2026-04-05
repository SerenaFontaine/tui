package widget

import (
	"time"
	"github.com/SerenaFontaine/tui"
)

// SpinnerStyle defines the frames of a spinner animation.
type SpinnerStyle struct {
	Frames   []string
	Interval time.Duration
}

// Predefined spinner styles.
var (
	SpinnerDots = SpinnerStyle{
		Frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		Interval: 80 * time.Millisecond,
	}
	SpinnerLine = SpinnerStyle{
		Frames:   []string{"|", "/", "-", "\\"},
		Interval: 100 * time.Millisecond,
	}
	SpinnerCircle = SpinnerStyle{
		Frames:   []string{"◐", "◓", "◑", "◒"},
		Interval: 120 * time.Millisecond,
	}
	SpinnerBounce = SpinnerStyle{
		Frames:   []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"},
		Interval: 100 * time.Millisecond,
	}
	SpinnerMeter = SpinnerStyle{
		Frames:   []string{"▱▱▱", "▰▱▱", "▰▰▱", "▰▰▰", "▰▰▱", "▰▱▱"},
		Interval: 150 * time.Millisecond,
	}
	SpinnerGlobe = SpinnerStyle{
		Frames:   []string{"🌍", "🌎", "🌏"},
		Interval: 200 * time.Millisecond,
	}
	SpinnerBlock = SpinnerStyle{
		Frames:   []string{"█", "▓", "▒", "░", "▒", "▓"},
		Interval: 100 * time.Millisecond,
	}
)

// SpinnerTickMsg triggers a spinner frame advance.
type SpinnerTickMsg struct {
	Time time.Time
	ID   int
}

// Spinner is an animated loading indicator.
type Spinner struct {
	SpinnerStyle SpinnerStyle
	Style        tui.Style
	Label        string
	frame        int
	id           int
}

var spinnerIDCounter int

// NewSpinner creates a new spinner with the dots style.
func NewSpinner() *Spinner {
	spinnerIDCounter++
	return &Spinner{
		SpinnerStyle: SpinnerDots,
		id:           spinnerIDCounter,
	}
}

// SetStyle sets the spinner text style.
func (s *Spinner) SetStyle(st tui.Style) *Spinner { s.Style = st; return s }

// SetSpinnerStyle changes the spinner animation.
func (s *Spinner) SetSpinnerStyle(ss SpinnerStyle) *Spinner { s.SpinnerStyle = ss; return s }

// SetLabel sets the label shown after the spinner.
func (s *Spinner) SetLabel(label string) *Spinner { s.Label = label; return s }

// Tick returns a Cmd that sends a SpinnerTickMsg after the interval.
func (s *Spinner) Tick() tui.Cmd {
	interval := s.SpinnerStyle.Interval
	id := s.id
	return func() tui.Msg {
		time.Sleep(interval)
		return SpinnerTickMsg{Time: time.Now(), ID: id}
	}
}

// Update handles tick messages.
func (s *Spinner) Update(msg tui.Msg) (*Spinner, tui.Cmd) {
	switch msg := msg.(type) {
	case SpinnerTickMsg:
		if msg.ID != s.id {
			return s, nil
		}
		s.frame = (s.frame + 1) % len(s.SpinnerStyle.Frames)
		return s, s.Tick()
	}
	return s, nil
}

// View returns the current spinner frame as a string.
func (s *Spinner) View() string {
	if len(s.SpinnerStyle.Frames) == 0 {
		return ""
	}
	return s.SpinnerStyle.Frames[s.frame]
}

// Render draws the spinner and label.
func (s *Spinner) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}
	frame := s.View()
	x := buf.SetString(area.X, area.Y, frame, s.Style)
	if s.Label != "" {
		buf.SetString(area.X+x, area.Y, " "+s.Label, s.Style)
	}
}
