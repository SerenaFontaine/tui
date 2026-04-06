package tui

import (
	"bytes"
	"testing"
	"time"
)

func TestScreenQueryAndReadTimeout(t *testing.T) {
	r := &slowReader{}
	var out bytes.Buffer
	s := &Screen{
		in:      r,
		out:     &out,
		managed: false,
		fd:      -1,
		events:  make(chan Msg, 64),
		done:    make(chan struct{}),
	}

	start := time.Now()
	_, err := s.QueryAndRead("\x1b_Gi=1;", 100*time.Millisecond)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("should timeout with error")
	}
	if elapsed < 50*time.Millisecond {
		t.Error("should wait near the timeout duration")
	}
	if elapsed > 500*time.Millisecond {
		t.Error("should not wait too long past timeout")
	}
}

type slowReader struct{}

func (r *slowReader) Read(p []byte) (int, error) {
	time.Sleep(10 * time.Second)
	return 0, nil
}

type nopWriter struct{}

func (w *nopWriter) Write(p []byte) (int, error) { return len(p), nil }
