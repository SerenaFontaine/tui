package tui

import "testing"

func TestImagePlacementWithAlt(t *testing.T) {
	p := NewImagePlacement(0, 0, 10, 5).WithAlt("Logo")
	if p.alt != "Logo" {
		t.Errorf("alt = %q, want 'Logo'", p.alt)
	}
}

func TestImagePlacementAltDefault(t *testing.T) {
	p := NewImagePlacement(0, 0, 10, 5)
	if p.alt != "" {
		t.Errorf("default alt = %q, want empty", p.alt)
	}
}

func TestImagePlacementWithAnimation(t *testing.T) {
	p := NewImagePlacement(0, 0, 10, 5).WithAnimation()
	if !p.IsAnimation() {
		t.Error("IsAnimation() should return true after WithAnimation()")
	}
}

func TestImagePlacementAnimationDefault(t *testing.T) {
	p := NewImagePlacement(0, 0, 10, 5)
	if p.IsAnimation() {
		t.Error("IsAnimation() should be false by default")
	}
}
