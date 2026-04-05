package tui

import "time"

// TickMsg is sent on each tick interval.
type TickMsg struct {
	Time time.Time
}

// TickCmd creates a command that sends a TickMsg after the given duration.
func TickCmd(d time.Duration) Cmd {
	return func() Msg {
		time.Sleep(d)
		return TickMsg{Time: time.Now()}
	}
}

// TickEvery is a convenience for TickCmd that takes a rate in FPS.
func TickEvery(fps float64) Cmd {
	d := time.Duration(float64(time.Second) / fps)
	return TickCmd(d)
}

// AfterCmd creates a command that sends a message after a delay.
func AfterCmd(d time.Duration, msg Msg) Cmd {
	return func() Msg {
		time.Sleep(d)
		return msg
	}
}

// PeriodicCmd creates a command that sends a message at regular intervals.
// The returned Cmd sends one message and should be re-invoked to continue.
func PeriodicCmd(interval time.Duration, fn func(time.Time) Msg) Cmd {
	return func() Msg {
		time.Sleep(interval)
		return fn(time.Now())
	}
}
