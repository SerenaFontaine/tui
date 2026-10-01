package tui

import (
	"bytes"
	"strings"
	"testing"
)

func TestScreenSuspendResumeMouseDisabled(t *testing.T) {
	var out bytes.Buffer
	s := &Screen{out: &out, fd: -1}
	s.Suspend()
	s.Resume()

	if strings.Contains(out.String(), "\x1b[?100") {
		t.Errorf("Suspend/Resume wrote mouse sequences with mouse disabled: %q", out.String())
	}
}

func TestScreenSuspendResumeMouseEnabled(t *testing.T) {
	var out bytes.Buffer
	s := &Screen{out: &out, fd: -1, mouse: true}
	s.Suspend()
	s.Resume()

	if !strings.Contains(out.String(), "\x1b[?1003l") || !strings.Contains(out.String(), "\x1b[?1003h") {
		t.Errorf("Suspend/Resume should toggle mouse tracking, got %q", out.String())
	}
}
