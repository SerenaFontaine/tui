package tui

import "testing"

func TestFocusManagerBasic(t *testing.T) {
	f := NewFocusManager("a", "b", "c")
	if f.Current() != "a" {
		t.Errorf("initial focus = %q, want 'a'", f.Current())
	}
	if f.Index() != 0 {
		t.Errorf("initial index = %d, want 0", f.Index())
	}
}

func TestFocusManagerNext(t *testing.T) {
	f := NewFocusManager("a", "b", "c")
	f.Next()
	if f.Current() != "b" {
		t.Errorf("after Next() = %q, want 'b'", f.Current())
	}
	f.Next()
	f.Next()
	if f.Current() != "a" {
		t.Errorf("Next() should wrap, got %q, want 'a'", f.Current())
	}
}

func TestFocusManagerPrev(t *testing.T) {
	f := NewFocusManager("a", "b", "c")
	f.Prev()
	if f.Current() != "c" {
		t.Errorf("Prev() from start should wrap to 'c', got %q", f.Current())
	}
}

func TestFocusManagerFocus(t *testing.T) {
	f := NewFocusManager("a", "b", "c")
	if !f.Focus("b") {
		t.Error("Focus('b') should return true")
	}
	if f.Current() != "b" {
		t.Errorf("after Focus('b'), Current() = %q", f.Current())
	}
	if f.Focus("nonexistent") {
		t.Error("Focus('nonexistent') should return false")
	}
}

func TestFocusManagerIsFocused(t *testing.T) {
	f := NewFocusManager("sidebar", "content")
	if !f.IsFocused("sidebar") {
		t.Error("should be focused on sidebar initially")
	}
	if f.IsFocused("content") {
		t.Error("should not be focused on content initially")
	}
}

func TestFocusManagerCount(t *testing.T) {
	f := NewFocusManager("a", "b", "c")
	if f.Count() != 3 {
		t.Errorf("Count() = %d, want 3", f.Count())
	}
}

func TestFocusManagerEmpty(t *testing.T) {
	f := NewFocusManager()
	if f.Current() != "" {
		t.Errorf("empty manager Current() = %q, want empty", f.Current())
	}
	if f.Next() != "" {
		t.Error("empty manager Next() should return empty")
	}
}
